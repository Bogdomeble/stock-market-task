#!/bin/bash

#
# run.sh - alternative to the docker compose approach
#

export PORT=${1:-8080}

docker-compose up --build
