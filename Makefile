migrate:
	go run cmd/migrate/main.go

feature:
	cd fujilane && godog

server:
	go run cmd/server/main.go
