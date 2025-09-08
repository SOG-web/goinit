.PHONY: build test clean release install

# Build the CLI tool
build:
	go build -o bin/goinit .

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/goinit-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o bin/goinit-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/goinit-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o bin/goinit-windows-amd64.exe .

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f goinit
	rm -f goinit.exe

# Install locally
install: build
	cp bin/goinit $(HOME)/.local/bin/goinit

# Create release archives
release: build-all
	mkdir -p release
	tar -czf release/goinit-linux-amd64.tar.gz -C bin goinit-linux-amd64
	tar -czf release/goinit-darwin-amd64.tar.gz -C bin goinit-darwin-amd64
	tar -czf release/goinit-darwin-arm64.tar.gz -C bin goinit-darwin-arm64
	zip release/goinit-windows-amd64.zip bin/goinit-windows-amd64.exe

# Development setup
dev-setup:
	go mod tidy
	go mod download

# Lint code
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...