# Makefile for epedia web scraper
build: clean ensure test
	@echo "Building..."
	go build
	go install
	@echo "Done!"

test: clean ensure
	@echo "Running tests..."
	go test ./...

ensure:
	@echo "Ensuring go dependencies exist..."
	dep ensure

clean:
	@echo "Cleaning binaries..."
	go clean