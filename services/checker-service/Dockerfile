FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY services/checker-service/go.mod ./services/checker-service/

COPY pkg/go.mod pkg/go.sum ./pkg/

WORKDIR /app/services/checker-service
RUN go mod download

WORKDIR /app
COPY . .

RUN go build -o /checker-app ./services/checker-service/cmd/app

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /checker-app .
CMD ["./checker-app"]