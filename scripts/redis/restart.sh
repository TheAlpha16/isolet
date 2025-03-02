#!/bin/bash

if docker ps -a --format '{{.Names}}' | grep -q '^isolet-redis$'; then
    docker rm -f isolet-redis
fi

docker run -d --name isolet-redis \
    -p 6379:6379 \
    --restart unless-stopped \
    redis:7 \
    redis-server --save "" --appendonly no
