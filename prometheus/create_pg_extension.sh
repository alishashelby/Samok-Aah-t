#!/bin/bash

psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" <<EOSQL
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
EOSQL