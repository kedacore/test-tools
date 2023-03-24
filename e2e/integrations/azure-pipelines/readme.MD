# New Pipelines for integration tests

The new demand-pipeline.yml is basically the same, but `demands` `golang` as this is the main requirement to proc the scaler.

new env vars:
AZURE_DEVOPS_DEMAND_PARENT_BUILD_DEFINITION_ID= (new id of build containing the demands from AzDO)