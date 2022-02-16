DOCKER_USER=saraceni
THIS_FILE := $(lastword $(MAKEFILE_LIST))
.PHONY: metalint test postgres-start postgres-stop dev dev-build start stop push-container pull-container update-docs

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
	go run cmd/main.go
dev-build:
	go build -o build/app cmd/main.go
	./build/app

# with docker-compose
start:
	docker-compose up -d
	docker-compose logs
stop:
	docker-compose down

# build and push image to registry (docker hub)
push-container:
	docker login -u $(DOCKER_USER)
	docker build -t $(DOCKER_USER)/suricate-bank .
	docker push $(DOCKER_USER)/suricate-bank
pull-container:
	docker pull saraceni/suricate-bank

# Update swagger docs
update-docs:
	swag init -g cmd/main.go