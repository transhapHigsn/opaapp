# syntax = docker/dockerfile:1

FROM golang:1.16-alpine AS base
WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM base AS build
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo -o /out/opaapp .

FROM alpine:latest
COPY --from=build /out/opaapp /app/opaapp
WORKDIR /app
CMD ["/app/opaapp"]
