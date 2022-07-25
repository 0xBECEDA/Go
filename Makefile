all: start_infra migrate

start_infra:
	docker-compose up -d

migrate:
	go run ./migrations/migrate.go
