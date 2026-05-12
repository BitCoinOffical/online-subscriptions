.PHONY: run docker-up test auth-test sub-test clean

APP_AUTH=./auth-service
APP_SUB=./subscription-service

run:
	go run cmd/main.go

docker-up:
	docker compose up --build

test: clean auth-test sub-test

auth-test:
	go -C $(APP_AUTH) test ./... -count=1

sub-test:
	go -C $(APP_SUB) test ./... -count=1

clean:
	go clean -testcache