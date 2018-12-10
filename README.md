# Manage SmartThings Event Notifications
This includes both a SmartThings Smart App that subscribes to user-defined events
and a Web service (in Go) to store and manage the events.

# Web Service Installation
* install mysql on the destination and only accept connections from localhost
* deploy the service code to the destination
* install Go on the destination

## Configuration
On the destination:
* edit the file `./config.yaml` and specify the DbDriver and DbDSN
* run `sh setup.sh` to create the database
