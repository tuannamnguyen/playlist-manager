#!/bin/bash

cd ../..
docker compose down db
docker volume prune --all --force
docker compose up -d db
