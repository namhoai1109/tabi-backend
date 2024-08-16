#!/usr/bin/env bash

# Ensure the database container is online and usable
echo "Waiting for database..."
until docker exec -i tabi-file.db psql -h localhost -U file -d file -c "SELECT 1" &> /dev/null
do
  # printf "."
  sleep 1
done
