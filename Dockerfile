FROM golang:1.27-alpine3.24 AS build

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN apk add --no-cache build-base sqlite

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o fooder-cli ./cmd/fooder-cli/main.go

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o fooder ./cmd/fooder/main.go

FROM alpine:3.24

RUN apk add --no-cache sqlite

WORKDIR /app
COPY --from=build /app/fooder .
COPY --from=build /app/fooder-cli .
COPY ./public /app/public/
COPY ./api/openapi /app/api/openapi/

EXPOSE 8080
CMD ["./fooder"]
