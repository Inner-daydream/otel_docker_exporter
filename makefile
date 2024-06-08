BINARY_NAME=otel-docker-exporter
VERSION=v1.0.0
PLATFORMS=darwin linux windows
ARCHITECTURES=amd64 386

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