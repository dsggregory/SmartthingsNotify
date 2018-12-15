package api

import (
	"code.dsg.com/smartthings_notif/stnotif"
	"github.com/gorilla/mux"
	"os"
	"testing"
)

var (
	f *stnotif.TestFixtures
	s server
)

func setupTest() error {
	f, err := stnotif.NewFixtures()
	//logrus.SetLevel(logrus.DebugLevel)
	if err == nil {
		s = server{appDir: "../../", config: f.Config, router: mux.NewRouter(), db: f.DbHandle}
		s.initRoutes()
	}
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
