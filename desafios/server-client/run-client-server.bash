#!/bin/bash
PWD=$(pwd)
BASENAME=$(basename $PWD)

if [ $BASENAME != "server-client" ]; then
    echo "Please, run this script from 'server-client' directory"
    exit 1
else
    # sudo docker compose up -d
    go run main.go
fi
