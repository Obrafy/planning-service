#!/bin/bash
set -e
cd "$(dirname $BASH_SOURCE)/.."

export PROJECT_NAME=${PROJECT_NAME:-`basename "$PWD"`}

if [ "$1" == "-h" ]; then
  echo "Project Name: $PROJECT_NAME - Usage: `basename $0` [...environments]"
  echo "[...environments] - A comma separated list of environments to run (e.g. alfa,production,development)"
  echo "If no environment is given, development is run"
  exit 0
fi

IFS=',' read -r -a environments <<< "$1"

# Check if no environment was given. If so, build development environment only
if ((${#environments[@]})); then
    # Iterate over environments and build images for each one
    for index in "${!environments[@]}"
    do
        echo "[${environments[index]}] Running Container"
        docker run -d --hostname $HOSTNAME --rm $PROJECT_NAME:${environments[index]} --environment ${environments[index]}
        echo "[${environments[index]}] Container Running Successfully"
    done  
else
    echo "[development] Running Container"
        docker run -d --hostname $HOSTNAME --rm $PROJECT_NAME:development --environment development
    echo "[development] Container Ruinning Successfully"
fi