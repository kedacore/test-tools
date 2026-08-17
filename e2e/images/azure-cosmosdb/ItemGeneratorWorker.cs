using Microsoft.Azure.Cosmos;

namespace CosmosDbTestTool
{
    /// <summary>
    /// Non-interactive replacement for the old interactive "ItemGenerator" console menu -
    /// inserts documents into the monitored container so change feed lag can be produced
    /// without a human at a terminal. Controlled entirely by environment variables so it
    /// can run headless as a Job/Pod.
    /// </summary>
    public class ItemGeneratorWorker : BackgroundService
    {
        private readonly ILogger<ItemGeneratorWorker> _logger;
        private readonly CosmosDbOptions _options;
        private readonly int _itemCount;
        private readonly int _intervalSeconds;

        public ItemGeneratorWorker(ILogger<ItemGeneratorWorker> logger, IConfiguration configuration)
        {
            _logger = logger;
            _options = CosmosDbOptions.FromEnvironment(configuration);
            _itemCount = int.TryParse(Environment.GetEnvironmentVariable("GENERATE_ITEM_COUNT"), out int count) ? count : 10;
            _intervalSeconds = int.TryParse(Environment.GetEnvironmentVariable("GENERATE_INTERVAL_SECONDS"), out int interval) ? interval : 0;
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            using var client = new CosmosClient(_options.Connection, new CosmosClientOptions { ConnectionMode = ConnectionMode.Gateway });
            Database database = await client.CreateDatabaseIfNotExistsAsync(_options.DatabaseId, cancellationToken: stoppingToken);
            Container container = await database.CreateContainerIfNotExistsAsync(_options.ContainerId, "/id", cancellationToken: stoppingToken);

            do
            {
                await GenerateItemsAsync(container, _itemCount, stoppingToken);

                if (_intervalSeconds <= 0)
                {
                    break;
                }

                await Task.Delay(TimeSpan.FromSeconds(_intervalSeconds), stoppingToken);
            } while (!stoppingToken.IsCancellationRequested);
        }

        private async Task GenerateItemsAsync(Container container, int itemsToInsert, CancellationToken cancellationToken)
        {
            _logger.LogInformation("Generating {Count} item(s) in {Database}/{Container}", itemsToInsert, _options.DatabaseId, _options.ContainerId);

            for (int i = 0; i < itemsToInsert && !cancellationToken.IsCancellationRequested; i++)
            {
                string id = Guid.NewGuid().ToString();
                await container.CreateItemAsync(
                    new { id, message = $"generated-item-{i}" },
                    new PartitionKey(id),
                    cancellationToken: cancellationToken);
            }

            _logger.LogInformation("Finished generating {Count} item(s)", itemsToInsert);
        }
    }
}
