build:
	@echo "Building backend..."
	@mkdir -p bin
	@go build -o bin/webserver cmd/api/main.go
	@chmod +x bin/webserver

run:
	@if [ -n "$$ENV" ]; then ENV=$$ENV go run cmd/api/main.go; else echo "ENV variable is not set"; fi

run-local:
	@echo "Running backend in local environment..."
	@ENV=local go run cmd/api/main.go

dev:
	@find . -type f | entr -r -s 'make run-local'

test:
	@find . -type f -name '*.go' | entr -r -s 'clear; go test ./...'