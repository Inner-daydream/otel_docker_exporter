BINARY_NAME=otel-docker-exporter
VERSION=1.2.1
PLATFORMS=darwin linux
ARCHITECTURES=amd64 arm64
IMAGE_TAG=ghcr.io/inner-daydream/otel_docker_exporter
build:
	@echo "Building $(BINARY_NAME) version $(VERSION)"
	@$(foreach platform,$(PLATFORMS), \
		$(foreach architecture,$(ARCHITECTURES), \
			GOOS=$(platform) GOARCH=$(architecture) go build -o './dist/$(BINARY_NAME)-$(VERSION)-$(platform)-$(architecture)' cmd/main.go; \
		) \
	)
	@echo "Build complete"

clean:
	@echo "Cleaning up"
	@rm -r ./dist

docker:
	@echo "Building docker image"
	@docker build -t $(IMAGE_TAG):$(VERSION) .
	@echo "Docker image built"
push: docker
	@echo "Pushing docker image"
	@docker push $(IMAGE_TAG):$(VERSION)
	@echo "Docker image pushed"