#!/bin/bash

if [[ ${#SSH_ORIGINAL_COMMAND[@]} != 1 ]]; then
    echo "Usage: update.sh <TAG>"
    echo "       update.sh 20250508.1315"
    exit 1
fi

TAG=${SSH_ORIGINAL_COMMAND[0]}
#TAG="latest"
IMAGE="lesheets"
#REMOTE_HOST="jtorr"
#REMOTE_DIR="/var/www/movements"
#COMPOSE_FILE="$REMOTE_DIR/docker-compose.yaml"
REMOTE_DIR="/root/services/lesheets"
REMOTE_HOST="ionos"
COMPOSE_FILE="$REMOTE_DIR/compose.yaml"

set -x

cd $REMOTE_DIR

# Load the image (from stdin) in the remote host
docker load

# Update the tag in the remote docker-compose file
sed -i "s|\(image:.*$IMAGE:\)[^[:space:]]*|\1$TAG|g" $COMPOSE_FILE
# force-recreate needed if it's the same tag
docker compose up -d
docker system prune -a -f
