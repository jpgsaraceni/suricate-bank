THIS_FILE := $(lastword $(MAKEFILE_LIST))
.PHONY: dev test start
dev:
	go run cmd/main.go
test:
	go clean -testcache
	go test ./...
start:
	docker-compose up -d
stop:
	docker-compose down