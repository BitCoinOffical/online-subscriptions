FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/server .

CMD [ "./server" ]