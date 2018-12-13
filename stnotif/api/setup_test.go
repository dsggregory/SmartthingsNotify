package api

import (
	"bufio"
	"code.dsg.com/smartthings_notif/stnotif/conf"
	"code.dsg.com/smartthings_notif/stnotif/dao"
	"encoding/csv"
	"io"
	"log"
	"os"
	"testing"
	"time"
)

var (
	s server
	c *conf.Conf
	d *dao.DbHandle
)

func loadFixtures(fpath string) {
	if d == nil {
		panic("dbHandle is nil")
	}
	csvFile, err := os.Open(fpath)
	if err != nil {
		log.Fatalf("%s: %+v", fpath, err)
	}
	defer csvFile.Close()

	nRecs := 0

	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("%s: %+v", fpath, err)
		}
		e := dao.NotifRec{}
		e.ID = 0
		t, err := time.Parse("01/02/2006 15:04:05", line[0])
		if err != nil {
			log.Printf("%v+: cannot import fixture data %+v\n", line, err)
			continue
		}
		e.EvTime = t.Unix()
		e.Device = line[1]
		e.Event = line[2]
		e.Value = line[3]
		e.Description = line[4]
		err = d.AddEvent(e)
		if err != nil {
			log.Println(err)
		} else {
			nRecs += 1
		}
	}
}

func setupTest() error {
	c = &conf.Conf{DbDriver: "mysql", DbDSN: "root@tcp(localhost)/test_st"}
	dbh, err := dao.NewDbHandlerTest(c)
	if err == nil {
		s = server{config: c, router: nil, db: dbh}
		d = dbh
	}

	loadFixtures("../../testdata/fixtures.csv")

	return err
}

func TestMain(m *testing.M) {
	err := setupTest()
	if err != nil {
		panic(err)
	} else {
		os.Exit(m.Run())
	}
}
