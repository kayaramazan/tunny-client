.PHONY: build install clean test

# Build the CLI
build:
	go build -o tunny ./cmd/tunny

# Install to $GOPATH/bin
install:
	go install ./cmd/tunny

# Clean build artifacts
clean:
	rm -f tunny
	rm -rf dist/

# Run tests
test:
	go test -v ./...

# Build for all platforms (using GoReleaser)
release-snapshot:
	goreleaser release --snapshot --clean --skip=publish

# Dev: build and run
dev:
	go run ./cmd/tunny connect localhost:3000 --dev

.DEFAULT_GOAL := build

