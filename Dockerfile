FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .
COPY config ./config
COPY templates ./templates
COPY migrations ./migrations

EXPOSE 9000

CMD ["./app", "--configFile", "./config/prod.yaml"]