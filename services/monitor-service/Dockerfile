FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY services/monitor-service/go.mod ./services/monitor-service/

COPY pkg/go.mod pkg/go.sum ./pkg/

WORKDIR /app/services/monitor-service
RUN go mod download

WORKDIR /app
COPY . .

RUN go build -o /monitor-app ./services/monitor-service/cmd/app

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /monitor-app .
CMD ["./monitor-app"]