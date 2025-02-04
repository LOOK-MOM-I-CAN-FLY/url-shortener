.PHONY: build run migrate test

build:
	docker-compose build

run:
	docker-compose up

migrate:
	docker-compose exec app goose -dir migrations postgres up

test:
	go test -v ./...