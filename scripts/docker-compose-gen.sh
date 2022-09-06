#!/bin/bash
set -e

# Vars 
REPLICAS="${1:-0}"
FILE_NAME="${2:-"docker-compose-dev.yaml"}"

BASE='
version: "3"
services:'

SERVER_BASE='
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - SERVER_PORT=12345
      - SERVER_LISTEN_BACKLOG=7
      - LOGGING_LEVEL=DEBUG
    networks:
      - testing_net'

BASE+="${SERVER_BASE}"

for (( i = 1; i <= $REPLICAS; i++ )) 
do
  
  CLIENT_BASE="
  client${i}:
    image: client:latest
    container_name: client${i}
    entrypoint: /client
    volumes:
      - ./.data/dataset-${i}.csv:/dataset.csv
    environment:
      - CLI_ID=${i}
      - CLI_SERVER_ADDRESS=server:12345
      - CLI_CONTESTANTS=./dataset.csv
    networks:
      - testing_net
    depends_on:
      - server"

  BASE+="${CLIENT_BASE}"
done

NETWORK_BASE="
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"

BASE+="${NETWORK_BASE}"

echo "${BASE}" > ${FILE_NAME}
