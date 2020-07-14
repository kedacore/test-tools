e2e-images: build-e2e-images push-e2e-images

build-e2e-images:
	./e2e/images/build.sh

push-e2e-images:
	./e2e/images/build.sh --push