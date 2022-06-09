#!/bin/bash

# Get platform
PLATFORM=$1
if [[ -z "$1" ]]; then 
    PLATFORM=linux
elif [[ "$1" == "mac" ]]; then 
    PLATFORM=darwin
fi

# Get architecture
ARCHITECTURE="$2"
if [[ -z "$2" ]]; then     
    ARCHITECTURE=amd64
fi

## Go to docker directory
DOCKER_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && cd ../docker && pwd )

# Get base Dockerfile
DOCKERFILE=$DOCKER_DIR/Dockerfile
NEW_DOCKERFILE=$DOCKER_DIR/Dockerfile-$PLATFORM-$ARCHITECTURE
TMP_DOCKERFILE=$DOCKER_DIR/Dockerfile.tmp

# Create new dockerfile if not existent
if [[ ! -f $NEW_DOCKERFILE ]]; then
    mkdir -p "${NEW_DOCKERFILE%/*}" && touch "$NEW_DOCKERFILE"
fi

# Create Dockerfile for platform and architecture
# awk '{sub("<PLATFORM>-<ARCHITECTURE>",'"$PLATFORM"')}1' $DOCKERFILE > $NEW_DOCKERFILE
sed "s/<PLATFORM>/$PLATFORM/g" $DOCKERFILE > $TMP_DOCKERFILE
sed "s/<ARCHITECTURE>/$ARCHITECTURE/g" $TMP_DOCKERFILE > $NEW_DOCKERFILE

rm $TMP_DOCKERFILE