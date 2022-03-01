# Suricate Bank

Suricate Bank is an api that creates accounts and transfers money between them. It is being built following Clean Arch.

A very special thanks to my Golang and Clean Arch mentor, [Helder](https://github.com/helder-jaspion). It's been a ride!

## Contents

* [Features](#features)
* [Coming Soon](#coming-soon)
* [Dependencies](#dependencies)
* [Getting Started](#getting-started)
  * [Running the App](#running-the-app)
  * [Automated Tests](#automated-tests-requires-docker-for-integration-tests)
  * [Manual Testing](#manual-testing)
  * [Available Routes](#available-routes)
    * [POST /accounts](#post-accounts)
    * [GET /accounts](#get-accounts)
    * [POST /login](#post-login)
    * [POST /transfers](#post-transfers)
    * [GET /transfers](#get-transfers)

## Features

* RESTful API
* Persistence on PostgreSQL DB with migrations
* Support for idempotent requests on create routes (account and transfer)
* Bearer Token (JWT) Auth on private routes (create transfer)
* Clean Architecture
* Containerized (Docker and Docker-Compose)
* Meaningful unit and integration tests
* Swagger Documentation
* Structured logging
* CI with GitHub Actions

## Coming Soon

* Apache Kafka events

## Dependencies

* [uuid](https://github.com/google/uuid) - For generating and inspecting unique account and transfer UUIDs;
* [pgx](https://github.com/jackc/pgx) - For configuring and connecting a pool to a PostgreSQL database and running queries;
* [Dockertest](https://github.com/ory/dockertest) - For running integration tests on temporary database containers;
* [migrate](github.com/golang-migrate/migrate) - For running database migrations;
* [bcrypt](https://golang.org/x/crypto/bcrypt) - For hashing and comparing hashed secrets;
* [jwt](github.com/golang-jwt/jwt/v4) - For signing and verifying JSON Web Tokens for authenticatioin;
* [chi](github.com/go-chi/chi) - For routing, panic recovery and requestID tracing;
* [cleanenv](github.com/ilyakaznacheev/cleanenv) - For reading .env and loading env variables;
* [redigo](github.com/gomodule/redigo) - For connecting and running commands on Redis (for idempotent HTTP requests);
* [Swaggo](https://github.com/swaggo) - For generating maintainable Swagger docs and Swagger UI
* [zerolog](https://github.com/rs/zerolog) - For structured logging.

## Getting started

To run this app in a container, the only requirement is [Docker Compose](https://docs.docker.com/compose/install/).

To run without a container, you will need [Go](https://go.dev/doc/install), [PostgreSQL](https://www.postgresql.org/download/), [Redis](https://redis.io/topics/quickstart) (configured and running), and optionally [Docker](https://docs.docker.com/get-docker/) to run integration tests.

An image of the application is available on Docker Hub registry [on this repository](https://hub.docker.com/r/saraceni/suricate-bank). You can pull it and use the postgres and redis instances you prefer (directly on your machine or building a container) to run the app. To pull the image, just run `make pull-container`.

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

You can also use Swagger UI for manual testing, as explained [below](#available-routes).

### Available routes

You can check out all routes on Swagger UI. Just access `http://HOST_RUNNING_THIS_APP:PORT_ITS_LISTENING_ON/swagger`. If you are running on the default URL, `http://localhost:8080/swagger`

#### POST `/accounts`

Create new account

##### Create account request payload example

```http
Content-Type: application/json

{
    "name": "another client",
    "cpf": "488.569.610-08",
    "secret": "really-good-one"
}
```

#### GET `/accounts`

List all accounts

#### GET `/accounts/{account_id}/balance`

Get account balance

#### POST `/login`

Login

##### Login request payload example

```http
Content-Type: application/json

{
    "cpf":"22061446035",
    "secret":"can't-tell-you"
}
```

#### POST `/transfers`

Create new transfer (requires Bearer token)

##### Create transfer payload example

```http
Content-type: application/json
Authorization: Bearer <JWT>
{
    "account_destination_id": "438e4746-fb04-4339-bd09-6cba20561835",
    "amount": 500
}
```

#### GET `/transfers`

List all transfers
