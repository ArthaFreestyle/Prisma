#!/bin/bash
set -e

echo "Running Postgre DB Migration......."
migrate -database "postgres://artha:passwordku@localhost:5432/prisma_db?sslmode=disable" -path ../db/migrations_postgre down
migrate -database "postgres://artha:passwordku@localhost:5432/prisma_db?sslmode=disable" -path ../db/migrations_postgre up
echo "Done!"