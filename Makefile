DOCKER_USER=saraceni
THIS_FILE := $(lastword $(MAKEFILE_LIST))
.PHONY: metalint test db-start db-stop dev dev-build start stop down push-container pull-container docs

# run linter
metalint:
ifeq (, $(shell which $$(go env GOPATH)/bin/golangci-lint))
	@echo "==> installing golangci-lint"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.44.0
endif
	$$(go env GOPATH)/bin/golangci-lint run --fix --allow-parallel-runners -c ./.golangci.yml ./...

# run tests
test:
	go test -race -count=1 ./...

# run without docker-compose
db-start:
	systemctl start postgresql
	redis-server --daemonize yes 
db-stop:
	systemctl stop postgresql
	redis-cli shutdown
dev:
	cp .env.example .env
	go run cmd/api/main.go
dev-build:
	cp .env.example .env
	go build -o build/app cmd/api/main.go
	./build/app

# with docker-compose
start:
	docker-compose up -d
	docker-compose logs
stop:
	docker-compose stop
down:
	docker-compose down

# build and push image to registry (docker hub)
push-container:
	docker login -u $(DOCKER_USER)
	docker build -t $(DOCKER_USER)/suricate-bank .
	docker push $(DOCKER_USER)/suricate-bank
pull-container:
	docker pull saraceni/suricate-bank

# Update swagger docs
docs:
	swag init -g cmd/api/main.go

.PHONY proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		protos/checking/checking.proto