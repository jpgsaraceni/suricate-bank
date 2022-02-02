THIS_FILE := $(lastword $(MAKEFILE_LIST))
.PHONY: dev dev-build test postgres-start postgres-stop start stop

# without docker-compose
dev:
	go run cmd/main.go
dev-build:
	go build -o build/app cmd/main.go
	./build/app
test:
	go clean -testcache
	go test ./...
postgres-start:
	systemctl start postgresql 
postgres-stop:
	systemctl stop postgresql

# with docker-compose
start:
	docker-compose up -d
	docker image prune --filter label=stage=gobuilder -f
	docker-compose logs
stop:
	docker-compose down
