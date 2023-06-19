# How to implement GitHub e2e GHA Tests for KEDA

Assumes you have already completed readme.MD

## Actions workflows:

You will need to modify the .env file to include the following variables:

.env file is located in the kedacore/tests directory.

- `GH_GHA_WORKFLOW_ID=` This is the workflow id of the workflow to run for scaled gha jobs. (test-runner-gha-job.yaml)
- `GH_APP_ID=` This is the app id of the github app that is used to create the runners.
- `GH_INST_ID=` This is the installation id of the github app that is used to create the runners.
- `GH_APP_KEY=` This is the private key of the github app that is used to create the runners.

### Run this after the test-tools PR is merged to main to get the workflow id for the workflow to run for scaled gha job.
1. `curl https://api.github.com/repos/kedacore/test-tools/actions/workflows`

2. Create a new github app for this repo. You need the app id when you have created it; turn off webhooks when you set it up.
3. Create a private key for the app and save it to a file. You need the private key when you have created it, base64 it and add it to the above env variable
4. Install the app on the repo. You need the installation id when you have installed it (Its in the URL at the top)
5. Permissions - Actions Read Only, Admin Read & Write, Metadata Read Only (in repositories)