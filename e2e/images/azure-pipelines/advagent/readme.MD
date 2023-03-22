# Mods for new E2E azure pipelines

## AdvAgent
The advagent is based on the pipelines parent design as laid out in the keda docs. Basically a mod of the traps to prevent killing when it is a `template` agent in `start.sh`

Also, don't use the docker stuff, as we don't need it. Sidecar docker in k8s is bad :)