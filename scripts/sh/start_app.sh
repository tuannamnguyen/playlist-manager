#!/bin/bash

cd ../..
docker compose --profile $1 up -d
cd ./cmd/api
go run .
