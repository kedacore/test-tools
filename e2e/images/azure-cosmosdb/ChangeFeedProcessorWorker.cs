using Microsoft.Azure.Cosmos;

namespace CosmosDbTestTool
{
    /// <summary>
    /// Runs a real Cosmos DB change feed processor against the monitored container, creating
    /// authentic .NET SDK lease documents in the lease container. This is what the KEDA
    /// azure-cosmosdb scaler e2e test needs bootstrapped before it can measure change feed lag -
    /// the scaler only reads lease/change-feed state via REST, it never creates leases itself.
    /// </summary>
    public class ChangeFeedProcessorWorker : BackgroundService
    {
        private readonly ILogger<ChangeFeedProcessorWorker> _logger;
        private readonly CosmosDbOptions _options;
        private CosmosClient? _client;
        private CosmosClient? _leaseClient;
        private ChangeFeedProcessor? _processor;

        public ChangeFeedProcessorWorker(ILogger<ChangeFeedProcessorWorker> logger, IConfiguration configuration)
        {
            _logger = logger;
            _options = CosmosDbOptions.FromEnvironment(configuration);
        }

        public override async Task StartAsync(CancellationToken cancellationToken)
        {
            var clientOptions = new CosmosClientOptions { ConnectionMode = ConnectionMode.Gateway };
            _client = new CosmosClient(_options.Connection, clientOptions);
            _leaseClient = _options.LeaseConnection == _options.Connection
                ? _client
                : new CosmosClient(_options.LeaseConnection, clientOptions);

            Database database = await _client.CreateDatabaseIfNotExistsAsync(_options.DatabaseId);
            Container monitoredContainer = await database.CreateContainerIfNotExistsAsync(_options.ContainerId, "/id");

            Database leaseDatabase = _options.LeaseDatabaseId == _options.DatabaseId
                ? database
                : await _leaseClient.CreateDatabaseIfNotExistsAsync(_options.LeaseDatabaseId);
            Container leaseContainer = await leaseDatabase.CreateContainerIfNotExistsAsync(_options.LeaseContainerId, "/id");

            _processor = monitoredContainer
                .GetChangeFeedProcessorBuilder<dynamic>(_options.ProcessorName, HandleChangesAsync)
                .WithInstanceName(Environment.MachineName)
                .WithLeaseContainer(leaseContainer)
                .WithErrorNotification(HandleErrorAsync)
                .Build();

            await _processor.StartAsync();
            _logger.LogInformation(
                "Change feed processor '{ProcessorName}' started on {Database}/{Container}, leases in {LeaseDatabase}/{LeaseContainer}",
                _options.ProcessorName, _options.DatabaseId, _options.ContainerId, _options.LeaseDatabaseId, _options.LeaseContainerId);

            await base.StartAsync(cancellationToken);
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            // All the work happens in the SDK's own change feed pump; just wait for shutdown.
            try
            {
                await Task.Delay(Timeout.Infinite, stoppingToken);
            }
            catch (TaskCanceledException)
            {
                // expected on shutdown
            }
        }

        public override async Task StopAsync(CancellationToken cancellationToken)
        {
            if (_processor != null)
            {
                await _processor.StopAsync();
            }
            await base.StopAsync(cancellationToken);
        }

        private Task HandleChangesAsync(IReadOnlyCollection<dynamic> changes, CancellationToken cancellationToken)
        {
            _logger.LogInformation("Processed {Count} change(s) from the change feed", changes.Count);
            return Task.CompletedTask;
        }

        private Task HandleErrorAsync(string leaseToken, Exception exception)
        {
            _logger.LogError(exception, "Unhandled exception on lease {LeaseToken}", leaseToken);
            return Task.CompletedTask;
        }

        public override void Dispose()
        {
            _client?.Dispose();
            if (_leaseClient != _client)
            {
                _leaseClient?.Dispose();
            }
            base.Dispose();
        }
    }
}
