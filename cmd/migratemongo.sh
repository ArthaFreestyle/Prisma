#!/bin/bash
set -e

echo "Running MongoDB Migration......."
migrate -database "mongodb://localhost:27017/prisma_db" -path /home/artha/Documents/PRISMA/db/migrations_mongo down
migrate -database "mongodb://localhost:27017/prisma_db" -path /home/artha/Documents/PRISMA/db/migrations_mongo up
echo "Done!"