FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/bot ./cmd/bot/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/bot .
COPY --from=builder /app/config/default-config.yaml ./config/

EXPOSE 8081
CMD ["./bot"]
