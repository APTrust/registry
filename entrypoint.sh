#!/bin/sh
# Copy all executables from ./bin to /usr/local/bin
# cp ./bin/linux/nsqadmin /usr/local/bin/
# cp ./bin/linux/nsqd /usr/local/bin/
# cp ./bin/linux/nsqlookupd /usr/local/bin/
# cp ./bin/linux/redis-cli /usr/local/bin/
# cp ./bin/linux/redis-server /usr/local/bin/

# Wait for Postgres to be ready
until PGPASSWORD=$DB_PASSWORD pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER; do
  echo "Waiting for Postgres..."
  sleep 2
done

# Run schema and migration files
echo "Running schema.sql"
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f /docker-entrypoint-initdb.d/schema.sql
for file in /docker-entrypoint-initdb.d/migrations/*.sql; do
  echo "Running migration $file"
  PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$file"
done

# Execute the original command
exec "$@"

