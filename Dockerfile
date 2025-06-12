FROM postgres:16

RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    python3-venv \
    gcc \
    libpq-dev \
    python3-dev \
    netcat-openbsd \
    gosu \
    && rm -rf /var/lib/apt/lists/*

RUN python3 -m venv /opt/venv
RUN /opt/venv/bin/pip install --upgrade pip \
    && /opt/venv/bin/pip install "psycopg2-binary>=2.9.6" \
    && /opt/venv/bin/pip install "patroni[etcd3]"

COPY patroni/post_init.sh /usr/local/bin/
COPY prometheus/create_pg_extension.sh /usr/local/bin/
COPY roles/analytic.sh /usr/local/bin/
COPY patroni/start.sh /start.sh

RUN chmod 755 /usr/local/bin/post_init.sh \
    && chmod 755 /usr/local/bin/create_pg_extension.sh \
    && chmod 755 /usr/local/bin/analytic.sh \
    && chmod +x /start.sh

RUN mkdir -p /data/patroni
ENV PATH="/opt/venv/bin:$PATH"
ENV PGDATA=/data/patroni

ENTRYPOINT ["/start.sh"]