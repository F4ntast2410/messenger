# === ЭТАП 1: Сборка бинарника ===
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Сначала только файлы модулей — для кэширования слоя go mod download.
# go.work требует go.mod/go.sum ВСЕХ модулей из use, даже тех, что не нужны
# этому конкретному сервису для сборки.
COPY go.work go.work.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY apps/chat-service/go.mod apps/chat-service/go.sum ./apps/chat-service/
COPY apps/user-service/go.mod apps/user-service/go.sum ./apps/user-service/
COPY apps/ws-gateway/go.mod apps/ws-gateway/go.sum ./apps/ws-gateway/

WORKDIR /build/apps/chat-service
RUN go mod download

# Теперь копируем сам код
WORKDIR /build
COPY pkg/pubsub ./pkg/pubsub
COPY apps/chat-service ./apps/chat-service

WORKDIR /build/apps/chat-service
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# === ЭТАП 2: Финальный образ ===
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /build/apps/chat-service/app .
EXPOSE 8080
CMD ["./app"]