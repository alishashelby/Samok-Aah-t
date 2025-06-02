FROM postgres:16

COPY roles/analytic.sh /docker-entrypoint-initdb.d/analytic.sh
RUN chmod +x /docker-entrypoint-initdb.d/analytic.sh

RUN echo "shared_preload_libraries = 'pg_stat_statements'" >> /usr/share/postgresql/postgresql.conf.sample && \
    echo "pg_stat_statements.max = 10000" >> /usr/share/postgresql/postgresql.conf.sample && \
    echo "pg_stat_statements.track = all" >> /usr/share/postgresql/postgresql.conf.sample

COPY prometheus/create_pg_extension.sql /docker-entrypoint-initdb.d/

EXPOSE 5432