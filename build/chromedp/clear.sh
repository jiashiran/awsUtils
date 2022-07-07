#!/usr/bin/env bash

CONTAINER_IDS=`docker ps -a | grep -v 'CONTAINER ID' | awk -F " " {'print $1'}`
for id in $CONTAINER_IDS
do
  echo "CONTAINER_ID=$id"
  docker stop $id
  docker rm -f $id
done


IMAGE_IDS=`docker images | grep -v '<none>' | grep -v 'IMAGE'  | awk -F " " {'print $3'}`
for id in $IMAGE_IDS
do
  echo docker rmi $id
  docker rmi $id
done
