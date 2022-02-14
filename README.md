# ⚠️ WORK IN PROGRESS (DEVELOPMENT RELEASE) ⚠️
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fjpgsaraceni%2Fsuricate-bank.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fjpgsaraceni%2Fsuricate-bank?ref=badge_shield)


Suricate Bank is an api that creates accounts and transfers money between them. It is being built following Clean Arch.

A very special thanks to my Golang and Clean Arch mentor, [Helder](https://github.com/helder-jaspion). It's been a ride!

## TODO

* Logging
* RequestID tracing
* Idempotency (redis)
* Panic recovery
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

An image of the application is available on Docker Hub registry [on this repository](https://hub.docker.com/r/saraceni/suricate-bank). You can pull it and use the postgres instance you prefer (directly on your machine or building a container) to run the app. To pull the image, just run `make pull-container`.

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

### Available routes

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


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fjpgsaraceni%2Fsuricate-bank.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fjpgsaraceni%2Fsuricate-bank?ref=badge_large)