FROM golang:1.24 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o sso ./cmd/sso

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/sso .
COPY ./config ./config
EXPOSE 50051
CMD ["./sso", "--config=./config/local.yaml"]
