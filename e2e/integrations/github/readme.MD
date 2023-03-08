# How to implement GitHub e2e Tests for KEDA

## Actions workflows:

You will need to modify the .env file to include the following variables:

.env file is located in the kedacore/tests directory.

- `GH_AUTOMATIONS_PAT=` This needs to be generated for an owner of this repository.  This is the PAT that will be used to create the runners and run the tests.
- `GH_OWNER=` This is the owner of the repository.  In this case it is `kedacore`
- `GH_SCOPE=` This is the scope of the runner. In this case it is `repo`
- `GH_REPOS=` This is the repo to test against.  In this case it is `test-tools`

### Run this after the test-tools PR is merged to main to get the workflow id for the workflow to run for scaled jobs and scaled objects.
1. `curl https://api.github.com/repos/kedacore/test-tools/actions/workflows`

- `GH_WORKFLOW_ID=` This is the workflow id of the workflow to run for scaled jobs. (test-runner-job.yaml)
- `GH_SO_WORKFLOW_ID=` This is the workflow id of the workflow to run for scaled objects. (test-runner-so.yaml)

