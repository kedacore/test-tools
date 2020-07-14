#! /bin/bash
set -e

DIR=$(dirname "$0")
cd $DIR

function push {
    for IMAGE_NAME in $(find * -name Dockerfile -exec dirname {} \; | tr '/' '-')
    do
        docker push docker.io/kedacore/$IMAGE_NAME:latest
    done
}

function build {
    for IMAGE in $(find * -name Dockerfile)
    do
        IMAGE_NAME=$(dirname $IMAGE | tr '/' '-')
        pushd $(dirname $IMAGE)
        docker build -t docker.io/kedacore/$IMAGE_NAME:latest .
        popd
    done
}

if [ "$1" == "--push" ]
then
    push
else
    build
fi
