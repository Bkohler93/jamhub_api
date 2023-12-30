#!/bin/bash

./scripts/migratereset.sh
./scripts/migrateup.sh

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o jamhubapi cmd/* && ./jamhubapi