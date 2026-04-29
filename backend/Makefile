IMAGE_TAG = todo-app-next-go-api

tools-install: ## Install development tools
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.3
	go install github.com/air-verse/air@v1.61.5

gen: ## Generate code from openapi
	oapi-codegen --config=./internal/adapter/controller/config.yaml ./api/openapi.yaml

external-up: ## Up external containers
	pushd ./build/docker && docker compose up -d postgres && popd

external-down: ## Down external containers
	pushd ./build/docker && docker compose down && popd

run: ## Run app
	go mod tidy
	APP_ENV=development go run ./cmd/server/main.go

dev: ## Run app with hot reload via air
	air

run-prod: ## Run app in production
	go mod tidy
	APP_ENV=production go run ./cmd/server/main.go

build-prod: ## Build app for production
	go mod download
	go build -o bin ./cmd/server/main.go

build-linux:
	go mod download
	GOOS=linux GOARCH=amd64 go build -v -o bin_linux ./cmd/server/main.go

start-prod: build-prod ## Build and run in production
	APP_ENV=production ./bin

start-linux:
	APP_ENV=production ./bin_linux

docker-build: ## Build image
	docker build --tag $(IMAGE_TAG) -f ./build/docker/Dockerfile .

docker-run: ## Run docker
	docker run -p 8080:8080 -it $(IMAGE_TAG)

docker-compose-up: docker-build ## Run docker compose up
	pushd ./build/docker && docker compose up -d --wait postgres web swagger-ui && popd

docker-compose-down: ## Run docker compose down
	pushd ./build/docker && docker compose down && popd

test: ## Run tests
	go test ./...

build: ## Build all packages (compile check)
	go build ./...

fmt: ## Format code
	golangci-lint fmt

tidy: ## Run go mod tidy
	go mod tidy

lint: ## Run lint
	golangci-lint run

vet: ## Run go vet
	go vet ./...

prettier: vet lint

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'