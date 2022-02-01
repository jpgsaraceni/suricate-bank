THIS_FILE := $(lastword $(MAKEFILE_LIST))
.PHONY: dev build test start stop localdb
dev:
	go run cmd/main.go
build:
	go build -o build/app cmd/main.go
	./build/app
test:
	go clean -testcache
	go test ./...
start:
	docker-compose up --build
stop:
	docker-compose down
localdb:
	systemctl start postgresql 