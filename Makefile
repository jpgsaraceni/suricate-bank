DOCKER_USER=saraceni
THIS_FILE := $(lastword $(MAKEFILE_LIST))
.PHONY: test postgres-start postgres-stop dev dev-build start stop push-container pull-container

# run tests
test:
	go test -race -count=1 ./...

# run without docker-compose
postgres-start:
	systemctl start postgresql 
postgres-stop:
	systemctl stop postgresql
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