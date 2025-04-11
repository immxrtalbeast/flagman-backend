FROM golang:alpine AS builder

# sudo docker run -e CONFIG_PATH=/app/config/local.yaml -p 8080:8080  flagman-back
RUN apk add --no-cache
WORKDIR /app
COPY . .

RUN go mod tidy

RUN go build -ldflags="-s -w" -o /app/main ./cmd/main.go

FROM alpine:latest
RUN apk add --no-cache 
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env /app/
COPY --from=builder /app/config ./config

EXPOSE 8080

ENTRYPOINT ["/app/main"]
CMD ["--config=/app/config/local.yaml"]  # Используем абсолютный путь