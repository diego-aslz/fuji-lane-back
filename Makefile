dependencies:
	go get github.com/golang/dep/cmd/dep
	dep ensure -v

	go get github.com/DATA-DOG/godog/cmd/godog

migrate:
	go run cmd/migrate/main.go

test:
	cd test && godog

server:
	go run cmd/server/main.go
