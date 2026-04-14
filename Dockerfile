FROM golang:bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api-gateway ./cmd/api-gateway/main.go

FROM debian:bookworm-slim
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/api-gateway .
COPY --from=builder /app/logs ./logs
COPY --from=builder /app/config ./config

EXPOSE 8000

CMD ["./api-gateway"]
