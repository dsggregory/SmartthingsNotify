# SmartThings Event Notifications WebApp
A Web service (in Go) to store and manage the SmartThings events.

# Run
With docker:
```text
$ docker build -t stnotif .
$ docker run -d -p 8080:8080 -e ALLOW_HOSTS='172.17.0.1,127.0.0.1' stnotif:latest
```
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
