#!/bin/sh

# set app's config to reach the database
[ -n "$MYSQL_SOCKET" ] && sed -i "s|socket:.*|socket: '$MYSQL_SOCKET'|g" ./config.yaml

# configure the database
export MAX_ALLOWED_PACKET="1M"
mysql_install_db --user=mysql --datadir=${DB_DATA_PATH}

sed -i "s|max_allowed_packet[ ]*=.*|max_allowed_packet = ${MAX_ALLOWED_PACKET}|g" /etc/mysql/my.cnf

# Restrict to localhost
cat >> /etc/mysql/my.cnf << EOF
[mysqld]
skip-networking
EOF

(cd '/usr' ; nohup /usr/bin/mysqld_safe --datadir=${DB_DATA_PATH} --nowatch; sleep 5)

cat > /tmp/sql << EOF
GRANT ALL PRIVILEGES ON *.* TO db_user @'localhost' IDENTIFIED BY 'db_passwd';
GRANT ALL PRIVILEGES ON *.* TO db_user @'127.0.0.1' IDENTIFIED BY 'db_passwd';
DELETE FROM mysql.user WHERE User='';
DROP DATABASE test;
FLUSH PRIVILEGES;
EOF
cat /tmp/sql | mysql -u root
rm -f /tmp/sql

sh ./create-db.sh
