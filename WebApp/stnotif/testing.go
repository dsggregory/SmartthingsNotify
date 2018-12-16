package stnotif

import (
	"bufio"
	"code.dsg.com/smartthings_notif/stnotif/conf"
	"code.dsg.com/smartthings_notif/stnotif/dao"
	"encoding/csv"
	"io"
	"log"
	"os"
	"time"
)

type TestFixtures struct {
	DbHandle *dao.DbHandle
	Config   *conf.Conf
}

func (tf *TestFixtures) AddFixture() {
	e := dao.NotifRec{}
	e.ID = 0
	e.EvTime = time.Now().Unix() - 1 // so other tests can use time.Now()
	e.Device = "fixture"
	e.Event = "add"
	e.Value = ""
	e.Description = ""
	tf.DbHandle.AddEvent(e)
}

func (tf *TestFixtures) loadFixtures(fpath string) {
	if tf.DbHandle == nil {
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
		err = tf.DbHandle.AddEvent(e)
		if err != nil {
			log.Println(err)
		} else {
			nRecs += 1
		}
	}
}

// NewFixtures creates and load test fixtures
func NewFixtures() (*TestFixtures, error) {
	tf := TestFixtures{}
	var c conf.Conf
	tf.Config = c.GetConf("../../config.yaml")
	tf.Config.Database.Database = "_test_st"

	dbh, err := dao.NewDbHandlerTest(tf.Config)
	if err == nil {
		tf.DbHandle = dbh
		fdir := "./"
		for i := 0; i < 5; i++ {
			fpath := fdir + "testdata/fixtures.csv"
			if _, err := os.Stat(fpath); !os.IsNotExist(err) {
				tf.loadFixtures(fpath)
				break
			} else {
				fdir = "../" + fdir // up one dir
			}
		}
	}

	return &tf, err
}
