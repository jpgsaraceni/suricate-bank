FROM golang:1.17-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine

WORKDIR /

COPY --from=build /app/main /

EXPOSE 8080

CMD ["./main"]