# ⚠️ WORK IN PROGRESS (DEVELOPMENT RELEASE) ⚠️

Suricate Bank is an api that creates accounts and transfers money between them. It is being built following Clean Arch.

A very special thanks to my Golang and Clean Arch mentor, [Helder](https://github.com/helder-jaspion). It's been a ride!

## TODO

* Logging
* RequestID tracing
* Idempotency (redis)
* Panic recovery
* Docker multistage build
* Swagger
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

To run this app in a container, the only requirement is [Docker Compose](https://docs.docker.com/compose/install/).

To run without a container, you will need [Go](https://go.dev/doc/install), [PostgreSQL](https://www.postgresql.org/download/) (configured and running), and optionally [Docker](https://docs.docker.com/get-docker/) to run integration tests.

### Running the app

1. Clone the project and enter directory:

```shell
git clone https://github.com/jpgsaraceni/suricate-bank.git && cd suricate-bank
```

2. Run in docker container:

```shell
make start
```

or without docker (after preparation instructed above):

```shell
make build
```

### Automated tests (requires Docker for integration tests)

```shell
make test
```

### Manual testing

The file `/client.http` can be used to test all available routes. I suggest using [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) VS Code extension for this (or any other HTTP request service you prefer).

### Available routes (TODO: document payloads and status codes)

#### POST `/accounts`

Create new account

#### GET `/accounts`

List all accounts

#### GET `/accounts/{account_id}/balance`

Get account balance

#### POST `/login`

Login

#### POST `/transfers`

Create new transfer (requires Bearer token)

#### GET `/transfers`

List all transfers
