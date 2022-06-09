#!/bin/bash
set -e
cd "$(dirname $BASH_SOURCE)/.."

TOOLS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && cd ../tools && pwd )
DOCKER_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && cd ../docker && pwd )

export PROJECT_NAME=${PROJECT_NAME:-`basename "$PWD"`}

if [ "$1" == "-h" ]; then
  echo "Project Name: $PROJECT_NAME - Usage: `basename $0` [...environments]"
  echo "[...environments] - A comma separated list of environments to build (e.g. alfa,production,development)"
  echo "If no environment is given, development is built"
  exit 0
fi

# Get platform
PLATFORM=$2
if [[ -z "$2" ]]; then 
    PLATFORM=linux
elif [[ "$2" == "mac" ]]; then 
    PLATFORM=darwin
fi

# Get architecture
ARCHITECTURE="$3"
if [[ -z "$3" ]]; then 
    ARCHITECTURE=amd64    
fi

# Create Dockerfile for platform and architecture
echo "Creating Dockerfile for platform and architecture: $PLATFORM - $ARCHITECTURE"
$TOOLS_DIR/transpile_dockerfile.bash $PLATFORM $ARCHITECTURE

# Get Dockerfile for platform and architecture
DOCKERFILE=$DOCKER_DIR/Dockerfile-$PLATFORM-$ARCHITECTURE


# Get first parameter as a comma separated list of environments to run
IFS=',' read -r -a environments <<< "$1"

# # Check if no environment was given. If so, build development environment only
if ((${#environments[@]})); then
    # Iterate over environments and build images for each one
    for index in "${!environments[@]}"
    do
        echo "[${environments[index]}] Running Building Image Command"
        docker image build \
            -t $PROJECT_NAME:${environments[index]} \
            -f $DOCKERFILE .
        echo "[${environments[index]}] Image Built Successfully"
    done  
else
    echo "No Environment Given - Building Development Environment"
    echo "[development] Running Building Image Command"
    docker image build \
        -t $PROJECT_NAME:development \
        -f $DOCKERFILE .
    echo "[development] Image Built Successfully"
fi