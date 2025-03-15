build:
	go run ./cmd

run:
	go run ./cmd

lint:
	golangci-lint run -c .golangci.yml

unit-test:
	go test ./... -short

generate-mocks:
	mockgen -source=./internal/handler.go -destination=./internal/mock_handler.go -package=internal
	mockgen -source=./internal/service.go -destination=./internal/mock_service.go -package=internal
	mockgen -source=./internal/repository.go -destination=./internal/mock_repository.go -package=internal
