FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

ADD jamhubapi /usr/bin/jamhubapi
COPY .env .env

ENV DOTENV_PATH .env

CMD ["jamhubapi"]