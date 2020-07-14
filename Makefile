IMAGE_TAG := $(shell git describe --always --abbrev=7)
e2e-images: build-e2e-images push-e2e-images

build-e2e-images:
	IMAGE_TAG=$(IMAGE_TAG) ./e2e/images/build.sh

push-e2e-images:
	IMAGE_TAG=$(IMAGE_TAG) ./e2e/images/build.sh --push