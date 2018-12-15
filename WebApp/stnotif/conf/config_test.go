package conf

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	c Conf
)

func TestMain(m *testing.M) {
	c = Conf{
		Database:     DbConf{Database: "db", User: "user", Password: "pass", Host: "host", Port: 3306, Driver: "driver", Socket: "sock"},
		AllowedHosts: []string{"foo.domain.com", "127.0.0.1", "192.*", "128.100.*"},
	}

	os.Exit(m.Run())
}

func TestConf_AllowsHost(t *testing.T) {
	assert := assert.New(t)

	assert.True(c.AllowsHost("foo.domain.com"))
	assert.True(c.AllowsHost("127.0.0.1"))
	assert.False(c.AllowsHost("localhost"))
	assert.True(c.AllowsHost("192.1.1.1"))
	assert.True(c.AllowsHost("128.100.1.1"))
	assert.False(c.AllowsHost("128.101.1.1"))

	c.AllowedHosts = []string{}
	assert.True(c.AllowsHost("any.domain.com"), "should allow any peer when not restricted")
}

// "user:passwd@tcp(host:port)/database?options"
func TestConf_DbDSN(t *testing.T) {
	assert := assert.New(t)

	c2 := c
	assert.Equal("user:pass@unix(sock)/db", c2.DbDSN())
	c2.Database.Socket = ""
	assert.Equal("user:pass@tcp(host:3306)/db", c2.DbDSN())
	c2.Database.Port = 0
	assert.Equal("user:pass@tcp(host:3306)/db", c2.DbDSN())
	c2.Database.Password = ""
	assert.Equal("user@tcp(host:3306)/db", c2.DbDSN())
}

func TestConf_DbDriver(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("driver", c.DbDriver())
}
