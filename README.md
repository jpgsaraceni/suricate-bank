# ⚠️ WORK IN PROGRESS ⚠️

Suricate Bank is an api that creates accounts and transfers money between them. It is being built following Clean Arch.

A very special thanks to my Golang and Clean Arch mentor, [Helder](https://github.com/helder-jaspion). It's been a ride!

## Dependencies

* [uuid](https://github.com/google/uuid) - For generating and inspecting unique account and transfer UUIDs;
* [pgx](https://github.com/jackc/pgx) - For configuring and connecting a pool to a PostgreSQL database and running queries;
* [Dockertest](https://github.com/ory/dockertest) - For running integration tests on temporary database containers;
* [bcrypt](https://golang.org/x/crypto/bcrypt) - For hashing and comparing hashed secrets;
* [jwt](github.com/golang-jwt/jwt/v4) - For signing and verifying JSON Web Tokens for authenticatioin;
* [chi](github.com/go-chi/chi) - For routing
* [Logrus](https://github.com/sirupsen/logrus) - For logging. This library is used in the dockertest example. I haven't set up logging for the project, so I will decide later if this will actually be used.

## Testing the app

So far I've implemented the internal layers (entities and usecases), repositories and handlers (missing main), so you can't actually run the app. However, you can test all the packages that have been created so far. To do so, clone the project:

```shell
git clone https://github.com/jpgsaraceni/suricate-bank.git
```

In your local project directory, install dependencies:

```shell
go mod download
```

Then run all unit and integration tests (you will need [Docker](https://www.docker.com/) installed. I had some trouble with permissions on Linux, [this](https://stackoverflow.com/questions/48568172/docker-sock-permission-denied) fixed it for me):

```shell
go test ./...
```
