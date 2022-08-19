#!/bin/bash

if [ "$1" == "stop" ]; then
    # sudo docker rm -f $(sudo docker ps -a -q)
    sudo docker-compose kill
    exit 0
fi

if [ "$1" == "clean" ]; then
    sudo rm -rf docker/db
    sudo mkdir -p docker/db
    exit 0
fi

if [ "$1" == "prod" ]; then
    sudo docker-compose pull && sudo docker-compose up -d --force-recreate --no-build
    exit 0
fi

cd api
CGO_ENABLED=0 go build
if [ $? -ne 0 ]; then
    echo "Error while building api"
    exit 1
fi
cd ..

sudo docker-compose up -d --force-recreate --build --remove-orphans

# shows api logs
if [ "$1" == "logs" ]; then
    sudo docker-compose logs -f $2
fi