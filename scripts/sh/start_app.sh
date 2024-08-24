#!/bin/bash

if [ $# -eq 0 ];
then
  echo "$0: Missing profile argument"
  exit 1
elif [ $# -gt 1 ];
then
  echo "$0: Too many arguments: $@"
  exit 1
fi


cd ../..
docker compose --profile $1 up -d
cd ./cmd/api
go run .
