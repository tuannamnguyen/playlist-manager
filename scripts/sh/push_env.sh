#!/bin/bash

if [ $# -eq 0 ];
then
  echo "$0: Missing environment arguments"
  exit 1
elif [ $# -gt 1 ];
then
  echo "$0: Too many arguments: $@"
  exit 1
fi

cd ../../cmd/api

npx dotenv-vault@latest push $1
