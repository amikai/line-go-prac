# syntax = docker/dockerfile:1

FROM golang:1.19-alpine AS build
WORKDIR /app
ENV CGO_ENABLED=0

COPY go.* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux GOARCH=amd64 go build -o /bin/app/cmd cmd/main.go

FROM scratch
COPY ./migrations /bin/app/static/migrations
COPY --from=build /bin/app/cmd /
CMD ["/cmd"]
