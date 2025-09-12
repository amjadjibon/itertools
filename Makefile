all: lint test

test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests passed"

lint:
	@echo "Running linter..."
	@golangci-lint run ./... --timeout 5m --fix
	@echo "Linter passed"

VERSION := 1.0.0
release:
	@git tag -a v$(VERSION) -m "Release version: $(VERSION)"
	@git push origin v$(VERSION)
	@echo "Released version: $(VERSION)"
