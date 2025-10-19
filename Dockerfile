FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/server/main.go
RUN go build -o insider-message-sender ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/insider-message-sender .
COPY --from=builder /app/docs ./docs
COPY .env.example .env

EXPOSE 8080
CMD ["./insider-message-sender"]
