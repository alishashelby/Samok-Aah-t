FROM postgres:16

RUN apt-get update && \
    apt-get install -y cron gettext && \
    rm -rf /var/lib/apt/lists/*

COPY backup.sh /backup.sh
COPY crontab.sh /crontab.sh
RUN chmod +x /backup.sh /crontab.sh

CMD ["/crontab.sh"]