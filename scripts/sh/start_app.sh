#!/bin/bash

cd ../..
docker compose up -d db pgadmin4
cd ./cmd/api
go run .
