#!/bin/bash

docker compose up --build --abort-on-container-exit migration-tester
docker compose down -v