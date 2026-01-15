FROM golang:1.24-alpine AS builder

WORKDIR /app

# 1. Копируем go.mod именно этого сервиса
COPY services/checker-service/go.mod ./services/checker-service/

# 2. Копируем go.mod общего пакета pkg (так как сервис от него зависит)
COPY pkg/go.mod pkg/go.sum ./pkg/

# 3. Переходим в папку сервиса и скачиваем его зависимости
WORKDIR /app/services/checker-service
RUN go mod download

# 4. Возвращаемся в корень и копируем весь исходный код
WORKDIR /app
COPY . .

# 5. Собираем бинарник (путь от корня /app)
RUN go build -o /checker-app ./services/checker-service/cmd/app

# Финальный этап
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /checker-app .
CMD ["./checker-app"]