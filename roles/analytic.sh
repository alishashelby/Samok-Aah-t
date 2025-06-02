#!/bin/bash

echo "ANALYST_NAMES: $ANALYST_NAMES"

IFS=',' read -r -a users <<< "$ANALYST_NAMES"

psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" <<EOSQL
DO \$\$
BEGIN
    CREATE ROLE analytic;

    GRANT USAGE ON SCHEMA public TO analytic;
    GRANT SELECT ON ALL TABLES IN SCHEMA public TO analytic;
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO analytic;

    RAISE INFO 'Analytic-reader created';
END
\$\$;
EOSQL


for user in "${users[@]}"; do
  echo "Processing user: $user"
  psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" <<EOSQL
  DO \$\$
  BEGIN
      CREATE ROLE "$user" WITH LOGIN PASSWORD '${user}_123' INHERIT;
      GRANT analytic TO "$user";

      RAISE INFO 'User % created', '${user}';
  END
  \$\$;
EOSQL
done