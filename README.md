# Project Template

This is a general Go project template. You can use it as a start for your projects.

Here you should describe your project description.

## Dependencies

As a database for the project we use [Postgres 13](https://www.postgresql.org/).

For monitoring and error tracking we use [Sentry.io](https://sentry.io).

## Development tools installation

### Golang

Golang is a main language for the project. 

You can install it with instructions from official website: [https://go.dev/doc/install](https://go.dev/doc/install).

### Migration

We use migration tool to update scheme of the database. 
You can find instruction for installation here: [https://github.com/golang-migrate/migrate/tree/master/cmd/migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

### Mocking

For mocking we use [https://github.com/vektra/mockery](https://github.com/vektra/mockery). 

You can download it with this command:
```sh
go install github.com/vektra/mockery/v2@latest
```

## Service configuration

For the correct service work, you need to use environment variables or the `.env` file.

In order to use the file, download the [godotenv](https://github.com/joho/godotenv)
binary file. General command to run `.env` file:

```bash
godotenv -f .env [command]
```

## Environment Variables

| Variable                |      Default       | Description                                                                                  |
|-------------------------|:------------------:|----------------------------------------------------------------------------------------------|
| LOG_LEVEL               |       debug        | Log level for logger. Possible options: trace, debug, info, warning, error, fatal and panic. |
| SERVER_PORT             |        8080        | Port on which app will run.                                                                  |
| SERVER_READ_TIMEOUT     |        15s         | App read response time.                                                                      |
| SERVER_WRITE_TIMEOUT    |        15s         | App write response time.                                                                     |
| DB_USER                 |        root        | Postgres database user.                                                                      |
| DB_PASS                 |      password      | Postgres database password.                                                                  |
| DB_HOST                 |         db         | Postgres database host.                                                                      |
| DB_PORT                 |        5432        | Postgres database port.                                                                      |
| DB_NAME                 |        app         | Postgres database name.                                                                      |
| TEST_DB_USER            |     test_root      | Postgres database user.                                                                      |
| TEST_DB_PASS            |      password      | Postgres database password.                                                                  |
| TEST_DB_HOST            |         db         | Postgres database host.                                                                      |
| TEST_DB_PORT            |        5432        | Postgres database port.                                                                      |
| TEST_DB_NAME            |      test_app      | Postgres database name.                                                                      |
| SENTRY_DSN              |                    | Sentry DSN.                                                                                  |
| SENTRY_ENV              |      staging       | Sentry environment.                                                                          |

## Installation

### Manual

You can build the service with the following terminal command:

```bash
go build ./cmd/...
```

Before running the service set up the environment variables and run the app:

```bash
app
```

### Docker installation

Installation could be done with `docker` and `docker compose` tools. Set up the environment variables or use `.env`
file.

```bash
# run only current module
docker compose up --build

# load all submodules
git submodule update --init --recursive

# run module with all submodules
docker compose -f docker compose-full.yml up --build

docker
```

## Migrations

Migrations are not running automatically on the service startup, so you need
the [golang-migrate tool](https://github.com/golang-migrate/migrate) to run them.

Database pattern: `postgresql://DB_USER:DB_PASS@DB_HOST:DB_PORT/DB_NAME?sslmode=disable`

```bash
# create migration
migrate create -ext sql -dir ./migrations [migration name]

# manual migrations
migrate -database [database url] -path ./migrations up

# docker compose migrations
docker compose run --rm deployment migrate -database [database url] -path ./migrations up

# migrations using in docker compose
docker compose run --rm app migrate -database postgresql://postgres:postgres@db:5432/postgres?sslmode=disable -path ./migrations up
```

## Testing

### Manual

Unit-tests are using mocks generated by [mockery](https://github.com/vektra/mockery). Mocks generation commands are
embedded in-source using [go:generate](https://golang.org/cmd/go/#hdr-Generate_Go_files_by_processing_source)
comments, so before running tests, you need to launch code generation.

Set up corresponding environment variables for tests.

```bash
# Run code generation
go generate ./...

# Run only unit tests
go test ./... -short

# with .env file
godotenv -f .env go test ./... -short
```

Integration tests need `TEST_...` environment variables.

```bash
# run tests
go test ./internal/module

# with .env file
godotenv -f .env go test ./internal/module
```

### Docker testing

Unit tests are performed during the building phase.

Integration tests could be run with these commands:

```bash
# integration tests
docker compose run --rm app module.test
```

## Documentation

[Swagger UI](https://swagger.io/tools/swagger-ui/) is a graphical representation of swagger specifications. You can
access REST documentation with "[base_endpoint]/docs/".

You can check availability of the application with "[base_endpoint]/status" endpoint.

### Note
I suggest to move pkg package to a separate repository. 