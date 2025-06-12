#!/bin/bash

echo "creating database"
psql -U "$POSTGRES_USER" -d postgres -c "CREATE DATABASE \"$POSTGRES_DB\""

echo "running analytic.sh"
/usr/local/bin/analytic.sh