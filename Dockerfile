FROM golang:1.24-alpine as builder

# Устанавливаем нужные зависимости для CGO и SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# ВАЖНО: CGO_ENABLED=1
RUN CGO_ENABLED=1 go build -o sso ./cmd/sso

FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/sso .
COPY ./config ./config
CMD ["./sso", "--config=./config/local.yaml"]
