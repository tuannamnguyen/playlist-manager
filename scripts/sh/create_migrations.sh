#!/bin/bash

read -p "Enter migrate message: " MIGRATE_MSG

migrate create -ext sql -dir ../migration/ -seq $MIGRATE_MSG
