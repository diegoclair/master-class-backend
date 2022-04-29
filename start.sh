#!/bin/sh

set -e

echo "run db migrations"
source /app/app.env  # load app.env file as environment variable before run migrations
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@" # this command will take the command received from CMD in DOCKERFILE