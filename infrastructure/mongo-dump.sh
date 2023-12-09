#!/bin/bash

CONTAINER_NAME="mongodb"
docker exec -it $CONTAINER_NAME bash -c 'mongodump --username=root --password=example --out=/data/db/backup'
