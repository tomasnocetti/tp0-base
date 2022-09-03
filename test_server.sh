#!/bin/bash
# set -v

SERVICE_NAME="${1:-server}"
PORT="${2:-12345}"
STRING_TEST='this is a great server'

### COLORS
GREEN='\033[0;32m'
BIGreen='\033[1;92m'
BIReed='\033[1;91m'
NC='\033[0m' # No Color

echo "${GREEN}## BUILDING IMAGE ## ${NC}"

docker build ./netcat -t netcat-test

echo "${GREEN}## TESTING SERVER ## ${NC}"
var=$(echo "${STRING_TEST}" | \
    docker run -i --rm --network=testing_net netcat-test nc ${SERVICE_NAME} ${PORT})

if [[ "$var" == *"$STRING_TEST"* ]]; then
  echo "${BIGreen}Congrats! Your server it's there."
else
  echo "${BIReed}Ups! No correct response received. ${var}"
fi
