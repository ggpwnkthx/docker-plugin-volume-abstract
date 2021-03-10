#!/bin/sh
PLUGIN_NAME=docker-plugin-volume-abstract
docker build -t $PLUGIN_NAME .
CONTAINER_ID=$(docker create $PLUGIN_NAME true)
mkdir rootfs
docker export "$CONTAINER_ID" | tar -x -C rootfs
docker rm -vf "$CONTAINER_ID"
docker rmi $PLUGIN_NAME
docker plugin create $PLUGIN_NAME .
