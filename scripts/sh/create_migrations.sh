#!/bin/bash

read -p "Enter migrate message: " MIGRATE_MSG

migrate create -ext sql -dir ../migrations/ -seq $MIGRATE_MSG
