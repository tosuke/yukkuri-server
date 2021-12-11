TARGET=dev

.PHONY: build
build:
	DOCKER_BUILDKIT=1 docker build -t yukkuri-server --platform linux/amd64 --target $(TARGET) --secret id=AQTK1_URL .
