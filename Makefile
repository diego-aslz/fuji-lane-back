dependencies:
	go mod download

migrate:
	go run cmd/migrate/main.go

reset:
	go run cmd/reset/main.go

seed:
	go run cmd/seed/main.go

feature:
	cd test && godog

server:
	go run cmd/server/main.go
