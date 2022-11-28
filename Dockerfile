FROM golang:1.18.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o api ./cmd/api/main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /app/api /

EXPOSE 8080

CMD ["./main"]