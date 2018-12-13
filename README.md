# Manage SmartThings Event Notifications
This includes both a SmartThings Smart App that subscribes to user-defined events
and a Web service (in Go) to store and manage the events.

# Build
```text
$ make
```

# Web Service Installation
* install mysql/mariadb on the destination and only accept connections from localhost
* deploy the service code to the destination
* install Go on the destination
* run `make` from the service code deployment directory

## Configuration
On the destination:
* edit the file `./config.yaml` and specify the DbDriver and DbDSN
* run `sh setup.sh` to create the database
