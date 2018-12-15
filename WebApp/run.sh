#!/bin/sh

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
    if [ "$ALLOW_HOSTS" == "any" ]; then
        sed -i "s|hosts:.*|hosts: []|g" ./config.yaml
    else
        sed -i "s|hosts:.*|hosts: [$ALLOW_HOSTS]|g" ./config.yaml
    fi
fi

./smartthings_notif
