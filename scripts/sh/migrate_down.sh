#!/bin/bash

read -p "Enter PostgreSQL username: " PGUSER
read -s -p  "Enter PostgreSQL password: " PGPASSWORD
echo
read -p "Enter PostgreSQL Host: " PGHOST
read -p "Enter PostgreSQL Port: " PGPORT

migrate -path ../migrations/ -database "postgres://$PGUSER:$PGPASSWORD@$PGHOST:$PGPORT/playlist_manager?sslmode=disable" down
