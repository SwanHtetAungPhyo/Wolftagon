OUT=./bin/app

build:
	@echo "Building the binary..."
	@go build -o $(OUT) ./cmd/main.go

docker:
	@echo "docker compose up..."
	@docker compose up -d
run:
	@echo "running the application..."
	@go run ./cmd/main.go