-include .$(env).sh

migrate:
	go run cmd/migrate/main.go

feature:
	cd fujilane && godog
