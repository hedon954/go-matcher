#!/bin/bash

echo "Waiting for middlewares to be ready..."
sleep 3

echo "Starting applications..."
APP_NAME=go-matcher-http-api SERVER_CONFIG_PATH=/etc/server_conf_tmp.yml /build/go-matcher-http-api &

tail -f /dev/null
