all: lint test

test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests passed"

lint:
	@echo "Running linter..."
	@golangci-lint run ./... --timeout 5m --fix
	@echo "Linter passed"
