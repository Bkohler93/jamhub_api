# FROM --platform=linux/amd64 debian:stable-slim
FROM golang:1.21

RUN apt-get update && apt-get install -y ca-certificates

ADD jamhubapi /usr/bin/jamhubapi

CMD ["jamhubapi"]