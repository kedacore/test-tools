#! /bin/bash
set -e

DIR=$(dirname "$0")

if [ -z "$IMAGE_TAG" ]
then
    IMAGE_TAG=latest
fi

cd $DIR

if [ "$1" == "--push" ]
then
    for IMAGE_NAME in $(find * -name Dockerfile -exec dirname {} \; | tr '/' '-')
    do
        docker image push -a ghcr.io/kedacore/tests-$IMAGE_NAME
    done
else
    for IMAGE in $(find * -name Dockerfile)
    do
        IMAGE_NAME=$(dirname $IMAGE | tr '/' '-')
        pushd $(dirname $IMAGE)
        docker build -t ghcr.io/kedacore/tests-$IMAGE_NAME:$IMAGE_TAG -t ghcr.io/kedacore/tests-$IMAGE_NAME:latest .
        popd
    done
fi
