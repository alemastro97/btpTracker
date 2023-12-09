#!/bin/bash

CONTAINER_NAME="mongodb"
docker exec -it $CONTAINER_NAME bash -c 'mongorestore  --username=root --password=example  /data/db/backup'