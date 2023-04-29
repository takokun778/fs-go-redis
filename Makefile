include .secret.env
export

.PHONY: redis
redis:
	@docker compose --project-name fs-go-redis --file ./.docker/docker-compose.yaml up -d

.PHONY: run
run:
	@go run main.go
