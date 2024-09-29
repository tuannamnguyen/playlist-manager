#!/bin/bash

read -p "Enter PostgreSQL username: " PGUSER
read -s -p "Enter PostgreSQL password: " PGPASSWORD
echo

migrate -path ../migrations/ -database "postgres://$PGUSER:$PGPASSWORD@localhost:5432/playlist_manager?sslmode=disable" up
