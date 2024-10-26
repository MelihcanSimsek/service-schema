#!/bin/bash

if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
  WINPTY="winpty"
else
  WINPTY=""
fi

$WINPTY docker run --name postgresql -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -p 6432:5432 -d postgres:16
echo "Postgresql started..."

until docker exec postgresql pg_isready -U postgres; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

$WINPTY docker exec -i postgresql psql -U postgres -d postgres -c "CREATE DATABASE product_service"
echo "DATABASE product_service created"

$WINPTY docker exec -i postgresql psql -U postgres -d product_service -c "
create table if not exists products
(
  id bigserial not null primary key,
  name varchar(255) not null,
  price double precision not null,
  discount double precision,
  store varchar(255) not null
);
"
echo "Table products created"
