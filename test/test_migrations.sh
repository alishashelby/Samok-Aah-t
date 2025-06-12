#!/bin/bash

docker-compose --env-file ../.env up --build --abort-on-container-exit migration-tester
docker-compose --env-file ../.env down -v