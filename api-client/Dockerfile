FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o api-client ./cmd/api

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/api-client .

EXPOSE 8080

CMD ["./api-client"]
