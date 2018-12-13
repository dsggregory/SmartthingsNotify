package api

import (
	"code.dsg.com/smartthings_notif/stnotif/conf"
	"code.dsg.com/smartthings_notif/stnotif/dao"
	"os"
	"testing"
)

var (
	s server
	c *conf.Conf
	d *dao.DbHandle
)

func setupTest() {
	c = &conf.Conf{DbDriver: "mysql", DbDSN: "root@tcp(localhost)/test_st"}
	d, _ = dao.NewDbHandlerTest(c)
	s = server{config: c, router: nil, db: d}
}

func TestMain(m *testing.M) {
	setupTest()
	os.Exit(m.Run())
}
