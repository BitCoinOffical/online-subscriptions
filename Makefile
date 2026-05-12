#run
run:
	go run cmd/main.go

docker-compose:
	docker compose up --build

test:
	cd ./auth-service && go clean -testcache && go test ./... && cd ../subscription-service && go clean -testcache && go test ./...
