FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/ /app/cmd/
COPY internal/ /app/internal/
COPY scripts/ /app/scripts/
COPY sql/ /app/sql/