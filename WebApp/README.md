# SmartThings Event Notifications WebApp
A Web service (in Go) to store and manage SmartThings events.

SmartThings events are sent to this service from the `Simple Event Logger` SmartApp found at https://github.com/krlaframboise/SmartThings/tree/master/smartapps/krlaframboise/simple-event-logger.src. When configuring the SmartApp, instead of configuring a Google spreadsheet and entering the GoogleSheets URL, you specify the URL of this service as:
```html
http://{host}:{port}/gs
```
This endpoint mimics the Google Sheets endpoints that the SmartApp expects.

# Run
The Docker release is going to be the easiest to install. It creates an image with the database installed and the package built. It secures the database and configures it to only accept local connections.

First, install Docker on the destination host where you want the web service to run, then clone this repo to that host and run:
```text
$ docker build -t stnotif .
```

## Config via Environment
The following environment variables passed to `docker run` are used to set configuration for the app. They are used
to modify the `config.yaml` file in the app of the docker image that was just built.

**ALLOW_HOSTS**

Determine from where connections should be allowed to the web app. Certainly the SmartThings servers, 
plus from where you may be browsing. Use these hosts as comma-separated values of the `ALLOW_HOSTS` environment variable 
that you pass to docker when you run the service. 
`ALLOW_HOSTS` supports globbing (ex. "128.15.*") so you must escape square brackets `[]` in IPv6 addresses.

**TZ**

Specifies the timezone of your SmartThings hub. This is required in order to properly convert notification timestamps
that do not specify the timezone. This setting may be in an IANA timezone format such as `America/New_York`.

## Run the Container
```text
$ docker run -d --name stnotif -p 8080:8080 \
    -e ALLOW_HOSTS='172.17.0.1,127.0.0.1' \
    -e TZ='America/New_York' \
    --mount 'source=stnotif-mysql,target=/var/lib/mysql' stnotif:latest
```

You can now browse to `http://{dockerhost}:8080/`.

## Example Deployment Script
```text
set -x
git pull
docker build -t stnotif . || exit 1
docker stop stnotif && \
docker container prune --force && \
docker run -d --name stnotif -p 8080:8080 --mount 'source=stnotif-mysql,target=/var/lib/mysql' -e ALLOW_HOSTS='127.0.0.1,98.25.156.124,54.83.15*' -e TZ='America/New_York' stnotif:latest
```

# Build
The following information is useful when building for local testing.
```text
$ make
```

## Web Service Installation
* install mysql/mariadb on the destination and only accept connections from localhost
* deploy the service code to the destination
* install Go on the destination
* run `make` from the service code deployment directory
* run `./smartthings_notif` to start the service

## Configuration
On the destination:
* edit the file `./config.yaml` and specify the Database configuration
* run `sh create-db.sh` to create the database. You can populate it with fixtures data by running `go run util/importcsv.go stnotif/testdata/fixtures.csv`
