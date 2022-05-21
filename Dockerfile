FROM golang:1.18.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /app/main /

EXPOSE 8080

CMD ["./main"]