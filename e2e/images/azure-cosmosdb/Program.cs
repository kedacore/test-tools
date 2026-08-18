using CosmosDbTestTool;

IHost host = Host.CreateDefaultBuilder(args)
    .ConfigureServices(services =>
    {
        services.AddHostedService<ChangeFeedProcessorWorker>();
    })
    .Build();

await host.RunAsync();
