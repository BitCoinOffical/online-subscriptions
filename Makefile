#run
run:
	go run cmd/main.go

docker-compose:
	docker compose up --build

test:
	go test ./...