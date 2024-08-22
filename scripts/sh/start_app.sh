#!/bin/bash

cd ../..
docker compose --profile test_minimal up -d
cd ./cmd/api
go run .
