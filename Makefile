.PHONY: run build tidy

# Run the dev server (auto-reloads with 'air' if installed)
run:
	go run ./cmd/main.go

# Build binary
build:
	go build -o bin/inkwell ./cmd/main.go

# Download dependencies
tidy:
	go mod tidy

# Install air for hot-reload (optional)
install-air:
	go install github.com/air-verse/air@latest

# Run with hot-reload (requires air)
dev:
	air -c .air.toml
