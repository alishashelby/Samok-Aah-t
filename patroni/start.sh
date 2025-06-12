#!/bin/bash

cleanup() {
    echo "stopping Patroni"
    patronictl exit --force
    sleep 5
    exit 0
}

trap cleanup SIGTERM SIGINT

mkdir -p /data/patroni
chown -R postgres:postgres /data
chmod 700 /data/patroni

if [ -f /data/patroni/postgresql.conf ]; then
    echo "reloading PostgreSQL configuration"
    gosu postgres pg_ctl reload
fi

until nc -z etcd0 2379 || nc -z etcd1 2379 || nc -z etcd2 2379; do
  echo "waiting for etcd"
  sleep 2
done

exec gosu postgres patroni /etc/patroni/patroni.yml