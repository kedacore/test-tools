#! /bin/bash
set -euo pipefail

# not all images can be built as multiarch at the moment
# here is a list of images that must be multiarch for e2e tests to pass
declare -A build_as_multiarch=(
    ["apache-ab"]=true
    ["websockets"]=true
)

DIR=$(dirname "$0")

if [[ -z "${IMAGE_TAG:-}" ]]; then
    IMAGE_TAG=latest
fi

options=$(getopt -l "push,platform,purge-policy:" -o "p,x,r:" -- "$@")
eval set -- "$options"

PUSH=false
PLATFORM=""
PURGE_POLICY="push"
while true; do
  case "$1" in
    -p|--push) PUSH=true; shift ;;
    -x|--platform) PLATFORM="$2"; shift 2 ;;
    -r|--purge-policy) PURGE_POLICY="$2"; shift 2 ;;
    --) shift; break ;;
    *) echo "Invalid option: $1" >&2; exit 1 ;;
  esac
done

cd $DIR

if [[ "$PUSH" == true ]]; then
    for IMAGE in $(find * -name Dockerfile); do
        IMAGE_NAME=$(dirname $IMAGE | tr '/' '-')
        if [[ "$PLATFORM" != "" && "${build_as_multiarch[$IMAGE_NAME]:-false}" == true ]]; then
            echo "building and pushing $IMAGE_NAME from $IMAGE for $PLATFORM"
            image_dir=$(dirname $IMAGE)
            docker buildx build --push --platform "$PLATFORM" -t "ghcr.io/kedacore/tests-$IMAGE_NAME" ./$image_dir
        else
            echo "only pushing $IMAGE_NAME"
            docker image push -a ghcr.io/kedacore/tests-$IMAGE_NAME
            if [[ "$PURGE_POLICY" == "push" ]]; then
                docker image rm ghcr.io/kedacore/tests-$IMAGE_NAME:$IMAGE_TAG ghcr.io/kedacore/tests-$IMAGE_NAME:latest
            fi
        fi
    done
else
    for IMAGE in $(find * -name Dockerfile); do
        echo "Checking disk space before pushing images..."
        df -h
        echo "Checking Docker disk usage..."
        docker system df
        IMAGE_NAME=$(dirname $IMAGE | tr '/' '-')
        pushd $(dirname $IMAGE)
        docker build --label "org.opencontainers.image.source=https://github.com/kedacore/test-tools" -t ghcr.io/kedacore/tests-$IMAGE_NAME:$IMAGE_TAG -t ghcr.io/kedacore/tests-$IMAGE_NAME:latest .
        if [[ "$PURGE_POLICY" == "build" ]]; then
            docker image rm ghcr.io/kedacore/tests-$IMAGE_NAME:$IMAGE_TAG ghcr.io/kedacore/tests-$IMAGE_NAME:latest
        fi
        popd
    done
fi
