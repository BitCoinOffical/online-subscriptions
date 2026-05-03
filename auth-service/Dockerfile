FROM golang:1.26-alpine AS builder

WORKDIR /auth

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o auth-server cmd/main.go

FROM alpine

WORKDIR /auth

COPY --from=builder /auth/migrations ./migrations
COPY --from=builder /auth/auth-server .

CMD [ "./auth-server" ]