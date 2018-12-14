package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"code.dsg.com/smartthings_notif/stnotif/conf"
	"code.dsg.com/smartthings_notif/stnotif/dao"

	log "github.com/sirupsen/logrus"
)

func main() {
	csvFile, err := os.Open(os.Args[1])
	if err != nil {
		log.WithError(err).Fatal("Usage: go run importcsv.go <csvfile>")
	}

	var c conf.Conf
	c.GetConf("./config.yaml")
	dbh, err := dao.NewDbHandler(&c)
	if err != nil {
		log.WithError(err).Fatal("config")
	}

	nRecs := 0

	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			log.Info("done")
			break
		} else if error != nil {
			log.WithError(error).Fatal("read")
		}
		e := dao.NotifRec{}
		e.ID = 0
		t, _ := time.Parse("01/02/2006 15:04:05", line[0])
		e.EvTime = t.Unix()
		e.Device = line[1]
		e.Event = line[2]
		e.Value = line[3]
		e.Description = line[4]
		err := dbh.AddEvent(e)
		if err != nil {
			log.WithError(err).WithField("rec", e).Error()
		} else {
			nRecs += 1
		}
	}

	fmt.Printf("Wrote %d records\n", nRecs)
}
