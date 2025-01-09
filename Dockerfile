FROM golang:1.24rc1-alpine3.21

# Instalacja narzędzi, które mogą być potrzebne
#RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/

EXPOSE 8080

RUN go build -o fooder /app/cmd/fooder/main.go

CMD ["./fooder"]
