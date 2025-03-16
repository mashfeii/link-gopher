FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/scrapper ./cmd/scrapper/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/scrapper .
COPY --from=builder /app/config/default-config.yaml ./config/

EXPOSE 8080
CMD ["./scrapper"]
