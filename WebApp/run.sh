#!/bin/sh

set -x
cd `dirname $0`

(cd '/usr' ; /usr/bin/mysqld_safe --datadir=${DB_DATA_PATH} --nowatch)

echo "Waiting on mysql to start"
n=0
while [ ! -e /run/mysqld/mysqld.sock ] ; do
    n=`expr $n + 1`
    if [ $n -ge 15 ]; then
        echo "giving up"
        break
    else
        echo "."
    fi
    sleep 1
done

if [ -n "$ALLOW_HOSTS" ]; then
    sed -i "s|allowedHosts:.*|allowedHosts: [$ALLOW_HOSTS]|g" ./config.yaml
fi

if [ -n "$HUB_TIMEZONE" ]; then
    sed -i "s|hubTimezone:.*|hubTimezone: '$HUB_TIMEZONE'|g" ./config.yaml
fi

./smartthings_notif
