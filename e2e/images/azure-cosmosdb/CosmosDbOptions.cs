namespace CosmosDbTestTool
{
    /// <summary>
    /// Resolves Cosmos DB configuration from either the ".NET style" double-underscore
    /// environment variables (CosmosDbConfig__X, bound automatically by IConfiguration) or
    /// the plain upper-snake-case variables (COSMOS_X). Both conventions are set by the KEDA
    /// azure_cosmosdb e2e test's deployment template, so either one alone is enough to run.
    /// </summary>
    public class CosmosDbOptions
    {
        public string Connection { get; private set; } = string.Empty;
        public string LeaseConnection { get; private set; } = string.Empty;
        public string DatabaseId { get; private set; } = string.Empty;
        public string ContainerId { get; private set; } = string.Empty;
        public string LeaseDatabaseId { get; private set; } = string.Empty;
        public string LeaseContainerId { get; private set; } = string.Empty;
        public string ProcessorName { get; private set; } = string.Empty;

        public static CosmosDbOptions FromEnvironment(IConfiguration configuration)
        {
            string Require(string configKey, string legacyEnvVar)
            {
                string? value = configuration[$"CosmosDbConfig:{configKey}"];
                if (string.IsNullOrEmpty(value))
                {
                    value = Environment.GetEnvironmentVariable(legacyEnvVar);
                }
                if (string.IsNullOrEmpty(value))
                {
                    throw new InvalidOperationException(
                        $"Missing required configuration. Set either 'CosmosDbConfig__{configKey}' or '{legacyEnvVar}'.");
                }
                return value;
            }

            var options = new CosmosDbOptions
            {
                Connection = Require("Connection", "COSMOS_CONNECTION"),
                DatabaseId = Require("DatabaseId", "COSMOS_DATABASE_ID"),
                ContainerId = Require("ContainerId", "COSMOS_CONTAINER_ID"),
                LeaseDatabaseId = Require("LeaseDatabaseId", "COSMOS_LEASE_DATABASE_ID"),
                LeaseContainerId = Require("LeaseContainerId", "COSMOS_LEASE_CONTAINER_ID"),
                ProcessorName = Require("ProcessorName", "COSMOS_PROCESSOR_NAME"),
            };

            // LeaseConnection defaults to the primary Connection when not set separately,
            // matching the KEDA azure-cosmosdb scaler's own metadata defaulting behavior.
            string? leaseConnection = configuration["CosmosDbConfig:LeaseConnection"];
            options.LeaseConnection = string.IsNullOrEmpty(leaseConnection) ? options.Connection : leaseConnection;

            return options;
        }
    }
}
