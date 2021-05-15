# OpaApp

Go Fiber webserver with integration to OPA

Note: Integration with OPA is wip. 😅

## Run server

```bash
go mod download
export ENV_FILE=/path/to/.env/file

# run directly
go run *.go

# or you can use build command
go build # build binary
./opaapp # execute binary

```

Sample .env file:

```.env
OPAAPP_ENV=development

OPAAPP_DB=tracker
OPAAPP_DB_USER=postgres
OPAAPP_DB_PASSWORD=passwd
OPAAPP_DB_HOST=localhost
OPAAPP_DB_PORT=5432

OPAAPP_PORT=3000

```

## Postgres DB for testing

You can use following to create test database:

```bash
docker run -it -e POSTGRES_PASSWORD=passwd -e POSTGRES_USER=postgres -e POSTGRES_DB=tracker -p 5432:5432 postgres
```
