# syntax=docker/dockerfile:1.4
FROM golang:1.24rc1-alpine3.21 AS build

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=arm64

RUN apk add --no-cache build-base sqlite

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# magic buildkit cache to save compilation
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o fooder-cli ./cmd/fooder-cli/main.go

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o fooder ./cmd/fooder/main.go

FROM alpine:3.21

RUN apk add --no-cache sqlite

WORKDIR /app
COPY --from=build /app/fooder .
COPY --from=build /app/fooder-cli .
COPY ./public /app/public/

EXPOSE 8080
CMD ["./fooder"]
