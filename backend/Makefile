.PHONY: build run test clean migrate-up migrate-down seed

build:
	@echo "Building..."
	@go build -o bin/server ./cmd/server

run:
	@go run ./cmd/server

test:
	@go test -v ./...

clean:
	@rm -rf bin/

migrate-up:
	@migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	@migrate -path migrations -database "$(DATABASE_URL)" down

seed:
	@go run ./cmd/seed/main.go
