FROM golang:1.26-alpine AS builder

WORKDIR /subs

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o subs-server cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /subs/migrations ./migrations
COPY --from=builder /subs/subs-server .

CMD [ "./subs-server" ]