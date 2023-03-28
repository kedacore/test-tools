#!/bin/bash

export RUNNER_ALLOW_RUNASROOT=1

TOKEN=GH_TOKEN_HERE

mkdir template
curl -o actions-runner-linux-arm64-2.303.0.tar.gz -L https://github.com/actions/runner/releases/download/v2.303.0/actions-runner-linux-arm64-2.303.0.tar.gz
tar xzf ./actions-runner-linux-arm64-2.303.0.tar.gz -C template

# Regular runners
for i in {1..25}
do
    mkdir keda-arm64-$i
    cp -r template/* keda-arm64-$i
    cd keda-arm64-$i
    echo HOME=$(pwd) > .env
    ./config.sh --url https://github.com/kedacore --token $TOKEN --name keda-arm64-$i --replace --unattended

    ./svc.sh install
    ./svc.sh start
    cd ..
done

# e2e runners
for i in {1..1}
do
    mkdir keda-arm64-e2e-$i
    cp -r template/* keda-arm64-e2e-$i
    cd keda-arm64-e2e-$i
    echo HOME=$(pwd) > .env
    ./config.sh --url https://github.com/kedacore --token $TOKEN --name keda-arm64-e2e-$i --labels e2e --replace --unattended
    
    ./svc.sh install
    ./svc.sh start
    cd ..
done

# http add-on e2e runers
for i in {1..2}
do
    mkdir keda-arm64-http-add-on-$i
    cp -r template/* keda-arm64-http-add-on-$i
    cd keda-arm64-http-add-on-$i
    echo HOME=$(pwd) > .env
    ./config.sh --url https://github.com/kedacore --token $TOKEN --name keda-arm64-http-add-on-$i --labels http-add-on-e2e --replace --unattended

    ./svc.sh install
    ./svc.sh start
    cd ..
done
