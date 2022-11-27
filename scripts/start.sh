#!/bin/sh

docker run --name bazaar --restart=always -d \
    -p 18528:18528 \
    -v /etc/bazaar:/etc/bazaar \
    -v /home/bazaar:/home/bazaar \
    -v /data/bazaar:/data/bazaar \
    bazaar:1.0 ./bazaar --path /etc/bazaar/config.toml
