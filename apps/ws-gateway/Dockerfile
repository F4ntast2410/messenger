# === ЭТАП 1: Сборка бинарника ===
FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.work go.work.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY apps/ws-gateway/go.mod apps/ws-gateway/go.sum ./apps/ws-gateway/
COPY apps/chat-service/go.mod apps/chat-service/go.sum ./apps/chat-service/
COPY apps/user-service/go.mod apps/user-service/go.sum ./apps/user-service/

WORKDIR /build/apps/ws-gateway
RUN go mod download

WORKDIR /build
COPY pkg/pubsub ./pkg/pubsub
COPY apps/ws-gateway ./apps/ws-gateway

WORKDIR /build/apps/ws-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# === ЭТАП 2: Финальный образ ===
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /build/apps/ws-gateway/app .
EXPOSE 8080
CMD ["./app"]