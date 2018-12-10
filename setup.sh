#!/bin/sh

mysql -uroot << EOF
CREATE DATABASE IF NOT EXISTS smartthings;
CREATE TABLE IF NOT EXISTS smartthings.notifications (
id int primary key auto_increment,
device_name varchar(64),
time bigint unsigned,
event varchar(64),
value varchar(64),
description varchar(128)
) ENGINE=INNODB;
EOF
