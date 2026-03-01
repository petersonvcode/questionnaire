build:
	@echo "Building backend..."
	@mkdir -p bin
	@go build -o bin/webserver cmd/api/main.go
	@chmod +x bin/webserver

run:
	@go run cmd/api/main.go

dev:
	@find . -type f | entr -r -s 'make run'

test:
	@find . -type f -name '*.go' | entr -r -s 'clear; go test ./...'