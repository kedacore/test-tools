#! /bin/bash
set -euo pipefail

# not all images can be built as multiarch at the moment
# here is a list of images that must be multiarch for e2e tests to pass
build_as_multiarch=(
    "apache-ab"
    "cloudevents-http"
    "mysql"
    "nsq"
    "prometheus"
    "rabbitmq"
    "redis-cluster-lists"
    "redis-lists"
    "redis-sentinel-lists"
    "websockets"
)

# Helper function to check if an element is in an array
contains_element () {
  local e match="$1"
  shift
  for e; do [[ "$e" == "$match" ]] && return 0; done
  return 1
}

DIR=$(dirname "$0")

if [[ -z "${IMAGE_TAG:-}" ]]; then
    IMAGE_TAG=latest
fi

options=$(getopt -l "push,platform:" -o "p,x:" -- "$@")
eval set -- "$options"

PUSH=false
PLATFORM=""
while true; do
  case "$1" in
    -p|--push) PUSH=true; shift ;;
    -x|--platform) PLATFORM="$2"; shift 2 ;;
    --) shift; break ;;
    *) echo "Invalid option: $1" >&2; exit 1 ;;
  esac
done

cd $DIR

if [[ "$PUSH" == true ]]; then
    for IMAGE in $(find * -name Dockerfile); do
        IMAGE_NAME=$(dirname $IMAGE | tr '/' '-')
        if [[ "$PLATFORM" != "" ]] && contains_element "$IMAGE_NAME" "${build_as_multiarch[@]}"; then
            echo "building and pushing $IMAGE_NAME from $IMAGE for $PLATFORM"
            image_dir=$(dirname $IMAGE)
            docker buildx build --push --platform "$PLATFORM" -t "ghcr.io/kedacore/tests-$IMAGE_NAME" ./$image_dir
        else
            echo "only pushing $IMAGE_NAME"
            docker image push -a ghcr.io/kedacore/tests-$IMAGE_NAME
        fi
    done
else
    for IMAGE in $(find * -name Dockerfile); do
        IMAGE_NAME=$(dirname $IMAGE | tr '/' '-')
        pushd $(dirname $IMAGE)
        docker build --label "org.opencontainers.image.source=https://github.com/kedacore/test-tools" -t ghcr.io/kedacore/tests-$IMAGE_NAME:$IMAGE_TAG -t ghcr.io/kedacore/tests-$IMAGE_NAME:latest .
        popd
    done
fi
