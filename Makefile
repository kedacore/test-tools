KEDA_TOOLS_GO_VERSION = 1.20.5

IMAGE_REGISTRY ?= ghcr.io
IMAGE_REPO     ?= kedacore
IMAGE_KEDA_TOOLS = $(IMAGE_REGISTRY)/$(IMAGE_REPO)/keda-tools:$(KEDA_TOOLS_GO_VERSION)

BUILD_PLATFORMS ?= linux/amd64,linux/arm64

E2E_IMAGE_TAG := $(shell git describe --always --abbrev=7)

##################################################
# e2e tests images                               #
##################################################

e2e-images: build-e2e-images push-e2e-images

build-e2e-images:
	IMAGE_TAG=$(E2E_IMAGE_TAG) ./e2e/images/build.sh

push-e2e-images:
	IMAGE_TAG=$(E2E_IMAGE_TAG) ./e2e/images/build.sh --push

##################################################
# tools image                                    #
##################################################

build-keda-tools:
	docker build -f tools/Dockerfile -t $(IMAGE_KEDA_TOOLS) --build-arg GO_VERSION=$(KEDA_TOOLS_GO_VERSION) .

push-keda-tools:
	docker buildx build --push --platform=${BUILD_PLATFORMS} -f tools/Dockerfile -t ${IMAGE_KEDA_TOOLS} --build-arg GO_VERSION=$(KEDA_TOOLS_GO_VERSION) .
