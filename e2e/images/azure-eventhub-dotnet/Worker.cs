using Azure.Messaging.EventHubs;
using Azure.Messaging.EventHubs.Processor;
using Azure.Storage.Blobs;
using System.Diagnostics;
using System.Text;

namespace EventHubWorker
{
    public class Worker : BackgroundService
    {
        string ehubNamespaceConnectionString = Environment.GetEnvironmentVariable("EVENTHUB_CONNECTION_STRING");
        string blobStorageConnectionString = Environment.GetEnvironmentVariable("STORAGE_CONNECTION_STRING");
        string blobContainerName = Environment.GetEnvironmentVariable("CHECKPOINT_CONTAINER");
        string consumerGroup = Environment.GetEnvironmentVariable("EVENTHUB_CONSUMERGROUP");
        private EventProcessorClient _processor;
        private readonly ILogger<Worker> _logger;

        public Worker(ILogger<Worker> logger)
        {
            _logger = logger;
            var storageClient = new BlobContainerClient(blobStorageConnectionString, blobContainerName);
            _processor = new EventProcessorClient(storageClient, consumerGroup, ehubNamespaceConnectionString);
            _processor.ProcessEventAsync += ProcessEventHandler;
            _processor.ProcessErrorAsync += ProcessErrorHandler;
        }

        public override Task StartAsync(CancellationToken cancellationToken)
        {
            _processor.StartProcessingAsync();
            return base.StartAsync(cancellationToken);
        }

        public override Task StopAsync(CancellationToken cancellationToken)
        {
            _processor.StopProcessingAsync();
            return base.StopAsync(cancellationToken);
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken) {}

        async Task ProcessEventHandler(ProcessEventArgs eventArgs)
        {
            // Write the body of the event to the console window
            Console.WriteLine("\tReceived event: {0}", Encoding.UTF8.GetString(eventArgs.Data.Body.ToArray()));

            // Update checkpoint in the blob storage so that the app receives only new events the next time it's run
            await eventArgs.UpdateCheckpointAsync(eventArgs.CancellationToken);
        }

        Task ProcessErrorHandler(ProcessErrorEventArgs eventArgs)
        {
            // Write details about the error to the console window
            Console.WriteLine($"\tPartition '{eventArgs.PartitionId}': an unhandled exception was encountered. This was not expected to happen.");
            Console.WriteLine(eventArgs.Exception.Message);
            return Task.CompletedTask;
        }
    }
}



