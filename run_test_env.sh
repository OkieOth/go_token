#!/bin/bash

scriptPos=${0%/*}

COMPOSE_FILE=$scriptPos/docker-compose/test_env.yaml

if ! docker compose -f $COMPOSE_FILE up --build --abort-on-container-exit --exit-code-from test_runner; then
  echo "error while running tests in docker compose :-/"
  exit 1
else
  exit 0
fi
