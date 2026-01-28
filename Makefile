KEDA_TOOLS_GO_VERSION = 1.25.6
K6_VERSION = v0.47.0

IMAGE_REGISTRY ?= ghcr.io
IMAGE_REPO     ?= kedacore
IMAGE_KEDA_TOOLS ?= $(IMAGE_REGISTRY)/$(IMAGE_REPO)/keda-tools:$(KEDA_TOOLS_GO_VERSION)
IMAGE_KEDA_K6_RUNNER ?= $(IMAGE_REGISTRY)/$(IMAGE_REPO)/keda-k6-runner

BUILD_PLATFORMS ?= linux/amd64,linux/arm64,linux/s390x

IMAGE_TAG := $(shell git describe --always --abbrev=7)

##################################################
# e2e tests images                               #
##################################################

e2e-images: build-e2e-images push-e2e-images

build-e2e-images:
	IMAGE_TAG=$(IMAGE_TAG) ./e2e/images/build.sh

push-e2e-images:
	IMAGE_TAG=$(IMAGE_TAG) ./e2e/images/build.sh --push --platform ${BUILD_PLATFORMS}

##################################################
# tools image                                    #
##################################################

build-keda-tools:
	docker build -f tools/Dockerfile -t $(IMAGE_KEDA_TOOLS) --build-arg GO_VERSION=$(KEDA_TOOLS_GO_VERSION) .

push-keda-tools:
	docker buildx build --push --platform=${BUILD_PLATFORMS} -f tools/Dockerfile -t ${IMAGE_KEDA_TOOLS} --build-arg GO_VERSION=$(KEDA_TOOLS_GO_VERSION) .

##################################################
# k6-runner image                                #
##################################################

build-keda-k6-runner:
	docker build -f k6-runner/Dockerfile -t ${IMAGE_KEDA_K6_RUNNER}:$(K6_VERSION) --build-arg K6_VERSION=$(K6_VERSION) .

push-keda-k6-runner:
	docker buildx build --push --platform=${BUILD_PLATFORMS} \
	-f k6-runner/Dockerfile \
	-t ${IMAGE_KEDA_K6_RUNNER}:latest \
	-t ${IMAGE_KEDA_K6_RUNNER}:$(K6_VERSION) \
	-t ${IMAGE_KEDA_K6_RUNNER}:$(K6_VERSION)-$(IMAGE_TAG) \
	--build-arg K6_VERSION=$(K6_VERSION) .
