#!/bin/sh

cd "$(dirname $0)"

static="$(realpath "$PWD/../front-end")"

podman run --rm --name nginx \
    -v "$PWD/conf":"/etc/nginx/templates" \
    -v "$static":"/www/" \
    -e "NGINX_PORT=80" \
    -e "NGINX_HOST=marketplace.test" \
    --network="host" \
    docker.io/library/nginx

