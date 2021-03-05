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
        docker push ghcr.io/kedacore/tests-$IMAGE_NAME:$IMAGE_TAG
        docker push docker.io/kedacore/tests-$IMAGE_NAME:$IMAGE_TAG
    done
else
    for IMAGE in $(find * -name Dockerfile)
    do
        IMAGE_NAME=$(dirname $IMAGE | tr '/' '-')
        pushd $(dirname $IMAGE)
        docker build -t docker.io/kedacore/tests-$IMAGE_NAME:$IMAGE_TAG .
        docker tag docker.io/kedacore/tests-$IMAGE_NAME:$IMAGE_TAG ghcr.io/kedacore/tests-$IMAGE_NAME:$IMAGE_TAG
        popd
    done
fi
