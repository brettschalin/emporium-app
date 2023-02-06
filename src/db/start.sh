#!/bin/sh

cd "$(dirname $0)"

podman run --rm --name pg \
    -v "$PWD/init":"/docker-entrypoint-initdb.d" \
    -e POSTGRES_USER="emporium" \
    -e POSTGRES_PASSWORD="secretpassword" \
    -e POSTGRES_DB="emporium" \
    -p "5432:5432" \
    docker.io/library/postgres
