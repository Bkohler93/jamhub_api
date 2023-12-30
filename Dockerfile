# Requires environment variables when running container, see .env.example


# FROM --platform=linux/amd64 debian:stable-slim
FROM golang:1.21

# set working directory
WORKDIR /app

# copy source code
COPY . .

# get dependencies
RUN go get -d -v ./...

# install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# build application
RUN go build -o jamhubapi ./cmd/*

EXPOSE 8080

CMD ["./scripts/run_prod.sh"]