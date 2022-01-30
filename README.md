# ⚠️ WORK IN PROGRESS ⚠️

Suricate Bank is an api that creates accounts and transfers money between them. It is being built following Clean Arch.

A very special thanks to my Golang and Clean Arch mentor, [Helder](https://github.com/helder-jaspion). It's been a ride!

## TODO

* Logging
* Request tracing
* Idempotency (redis)
* Docker
* GitHub Actions

## Dependencies

* [uuid](https://github.com/google/uuid) - For generating and inspecting unique account and transfer UUIDs;
* [pgx](https://github.com/jackc/pgx) - For configuring and connecting a pool to a PostgreSQL database and running queries;
* [Dockertest](https://github.com/ory/dockertest) - For running integration tests on temporary database containers;
* [migrate](github.com/golang-migrate/migrate) - For running database migrations;
* [bcrypt](https://golang.org/x/crypto/bcrypt) - For hashing and comparing hashed secrets;
* [jwt](github.com/golang-jwt/jwt/v4) - For signing and verifying JSON Web Tokens for authenticatioin;
* [chi](github.com/go-chi/chi) - For routing;
* [cleanenv](github.com/ilyakaznacheev/cleanenv) - For reading .env and loading env variables;
* [Logrus](https://github.com/sirupsen/logrus) - For logging. This library is used in the dockertest example. I haven't set up logging for the project, so I will decide later if this will actually be used.

## Getting started

### Requirements (while I haven't containerized the app)

1. [Go](https://go.dev/dl/)
2. [Docker](https://docs.docker.com/get-docker/) (for integration tests). *If you're using Linux, you might have problems with permissions for running the integration tests (they create a docker image). [this](https://stackoverflow.com/questions/48568172/docker-sock-permission-denied) solved it for me.*
3. [Postgres](https://www.postgresql.org/download/) (run a database and create a .env file accordingly. View `/cmd/.env.example`.)

### Automated testing

```shell
go test ./...
```

### Running the app

1. Clone the project and enter directory:

```shell
git clone https://github.com/jpgsaraceni/suricate-bank.git && cd suricate-bank
```

2. Install dependencies:

```shell
go mod download
```

3. Run the app

```shell
go run cmd/main.go
```

### Manual testing

The file `/client.http` can be used to test all available routes. I suggest using [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) VS Code extension for this (or any other HTTP request service you prefer).

### Available routes (TODO)
