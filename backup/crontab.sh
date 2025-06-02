#!/bin/bash

echo "[$(date)] Starting cron configuration"

printenv | grep -E 'POSTGRES_|^BACKUP_' > /etc/environment
echo "PATH=$PATH:/usr/bin:/usr/local/bin" >> /etc/environment

echo "SHELL=/bin/bash" > /etc/cron.d/backup
echo "BASH_ENV=/etc/environment" >> /etc/cron.d/backup
echo "$BACKUP_INTERVAL_CRON root /backup.sh >> /var/log/cron.log 2>&1" >> /etc/cron.d/backup

chmod 0644 /etc/cron.d/backup
cron
tail -f /dev/null