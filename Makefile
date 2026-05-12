#run
run:
	go run cmd/main.go

docker-compose:
	docker compose up --build

test:
	cd ./auth-service && go test ./... && cd ../subscription-service && go test ./...
