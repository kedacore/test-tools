using CosmosDbTestTool;

// RUN_MODE selects which worker to host in this single combined image:
//   "processor" (default) - runs a real Cosmos DB change feed processor (creates lease docs)
//   "generate"            - inserts test documents into the monitored container, non-interactively
string runMode = Environment.GetEnvironmentVariable("RUN_MODE") ?? "processor";

IHost host = Host.CreateDefaultBuilder(args)
    .ConfigureServices(services =>
    {
        switch (runMode.ToLowerInvariant())
        {
            case "generate":
                services.AddHostedService<ItemGeneratorWorker>();
                break;
            case "processor":
                services.AddHostedService<ChangeFeedProcessorWorker>();
                break;
            default:
                throw new InvalidOperationException($"Unknown RUN_MODE '{runMode}'. Expected 'processor' or 'generate'.");
        }
    })
    .Build();

await host.RunAsync();
