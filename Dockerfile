# used newer version to my local one just out of curiosity
FROM golang:1.24rc1-alpine3.21 AS base

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=arm64

RUN apk add --no-cache build-base sqlite

FROM base AS build
WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/


RUN go build -o fooder /app/cmd/fooder/main.go

EXPOSE 8080
CMD ["./fooder"]
