# ⚠️ WORK IN PROGRESS ⚠️ 

Suricate Bank is an api that creates accounts and transfers money between them. It is being built following Clean Arch.

A very special thanks to my Golang and Clean Arch mentor, [Helder](https://github.com/helder-jaspion). It's been a ride!

## Dependencies

* [uuid](https://github.com/google/uuid) - For generating and inspecting unique account and transfer UUIDs;
* [pgx](https://github.com/jackc/pgx) - For configuring and connecting a pool to a PostgreSQL database and running queries;
* [Dockertest](https://github.com/ory/dockertest) - For running integration tests on temporary database containers;
* [bcrypt](https://golang.org/x/crypto/bcrypt) - For hashing and comparing hashed secrets;
* [Logrus](https://github.com/sirupsen/logrus) - For logging. This library is used in the dockertest example. I haven't set up logging for the project, so I will decide later if this will actually be used.
