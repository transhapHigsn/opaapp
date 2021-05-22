# OpaApp

Go Fiber webserver for learning different technologies

Note: I have gone from implementing a Rest API with Open Policy Agent (OPA) to implementing:

- OpenTelemetry using LightStep
- Open Policy Agent
- Postgres

I think I still haven't strayed far enough, but soon I will be. This is what i'm thinking to add in this project:

- OpenAPI
- Pub/Sub

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

P.S.: I might not use some of them in actual projects, but whatever.

## Use docker to setup everything

- Export environment variable for docker buildkit: `export DOCKER_BUILDKIT=1`
- Build image: `docker build --rm -t opaap -f Dockerfile .`
- Run image: `docker run -it --rm -env-file .env.sample -w /app --network host opaap /app/opaapp`

note: `--network host` is for cases where you want to connect database running locally.

## Use pack to setup everything

- Build pack: `pack build opa-app --path . --builder paketobuildpacks/builder:tiny`
- Run server: `docker run -it --rm --network host --env-file .env.sample opa-app opaapp`

## Issues & Notes

- I have tried using scratch image in multi-stage build, but lightstep exporter started throwing authentication failed due to invalid x.509 certificate error. Due to that, difference between final image ~5 MB (alpine is bigger).
- When running app using prefork mode in docker, the app crashes. I am aware that this is a[known issue](https://github.com/gofiber/fiber/issues/1021#issuecomment-730537971) but the suggested solution isn't working. I will do some more debugging to find the exact issue. This is relevant to both `docker`-only and `pack` based deployment.
