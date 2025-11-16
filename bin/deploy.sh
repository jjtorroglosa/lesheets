#!/usr/bin/env bash

if [[ -z "$1" ]]; then
    echo "Usage: deploy.sh --ask"
    echo "Usage: deploy.sh <TAG>"
    echo "       deploy.sh 20250508.1315"
    exit 1
fi

IMAGE_NAME="lesheets"
IMAGE="$IMAGE_NAME" # this would be the fqdn of the image, e.g. git.jjtorroglosa.com/jtorr/$IMAGE_NAME
REMOTE_HOST="${REMOTE_HOST:-ionos1}"
DEFAULT_KEY="$HOME/.ssh/lesheets_deployer"
SSH_KEY="${SSH_KEY:-$DEFAULT_KEY}"
SSH="ssh -i $SSH_KEY $REMOTE_HOST"

get_tags() {
    curl -n -s https://git.jjtorroglosa.com/v2/jtorr/$IMAGE_NAME/tags/list |jq -r '.tags[]'
}

if [[ "$1" = "--ask" ]]; then
    TAG=$(get_tags|grep -v latest|sort -r|fzf --no-sort)
    if [[ -z "$TAG" ]]; then
        echo "Tag wasn't selected"
        exit 1
    fi
else
    TAG=$1
fi

echo Deploying $IMAGE:"$TAG" to "$REMOTE_HOST"

set -x

# Pull the image locally
# docker pull $IMAGE:"$TAG"

# Load the image in the remote host
docker save $IMAGE:"$TAG" |gzip | $SSH "$TAG"
