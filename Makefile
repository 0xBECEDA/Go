all:
	go build server.go
	go build client.go

start_infra:
	docker-compose up -d

migrate:
	go run ./migrations/migrate.go