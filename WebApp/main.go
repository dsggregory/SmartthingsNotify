package main

import (
	"os"

	"code.dsg.com/smartthings_notif/stnotif/conf"

	"code.dsg.com/smartthings_notif/stnotif/api"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	log.SetReportCaller(false)
}

func main() {
	var c conf.Conf
	c.GetConf("./config.yaml")
	if c.Debug {
		log.SetLevel(log.DebugLevel)
	}
	pwd, _ := os.Getwd()
	log.WithField("pwd", pwd).Debugf("Config => %+v", c)
	api.StartServer(&c)
}
