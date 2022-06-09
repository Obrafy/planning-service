#!/bin/bash
source ./crosscompile.bash
echo $1 $2

OPTS=${OPTS:-""}
CURRENT_BRANCH=-$([[ -d .git ]] && git log -1 --pretty=format:'%h' || echo 'docker')

if [ $CURRENT_BRANCH == "-master" ]; then
	CURRENT_BRANCH=
fi

CURRENT_TAG=$([[ -d .git ]] && git describe --always --long --dirty || echo 'docker')
VERSION_ID=$CURRENT_TAG$CURRENT_BRANCH
PROJECT_DIR=`pwd`
BIN_DIR=$PROJECT_DIR"/bin"

echo $BIN_DIR
set -e

if [[ "$1" == "mac" || "$1" == "darwin" ]]; then
    echo "mac or darwin"
    go-darwin-amd64 build -ldflags="-X main.version=$VERSION_ID" -v -o $BIN_DIR/darwin-amd64/app $PROJECT_DIR; \
    mkdir -p $BIN_DIR/darwin-amd64/configuration/ && cp $PROJECT_DIR/configuration/*.yaml $BIN_DIR/darwin-amd64/configuration; \
elif [[ "$1" == "linux" && "$2" == "arm" ]]; then
    echo "linux arm"
    go-linux-arm build -ldflags="-X main.version=$VERSION_ID" -v -o $BIN_DIR/linux-arm/app $PROJECT_DIR; \
    mkdir -p $BIN_DIR/linux-arm/configuration/ && cp $PROJECT_DIR/configuration/*.yaml $BIN_DIR/linux-arm/configuration; \
else
    echo "linux amd64"
    go-linux-amd64 build -ldflags="-X main.version=$VERSION_ID" -v -o $BIN_DIR/linux-amd64/app $PROJECT_DIR; \
    mkdir -p $BIN_DIR/linux-amd64/configuration/ && cp $PROJECT_DIR/configuration/*.yaml $BIN_DIR/linux-amd64/configuration; \
fi

if [ $? -ne 0 ]; then
    echo "ERROR BUILDIN GO"
    exit 1
fi