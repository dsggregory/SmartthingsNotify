package conf

import (
	"fmt"
	"github.com/gobwas/glob"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// DbConf is the Database config options
type DbConf struct {
	Driver   string `yaml:"driver"`
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Socket   string `yaml:"socket"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Conf an instantiated configuration
type Conf struct {
	ServerPort   int      `yaml:"serverPort"`
	Debug        bool     `yaml:"debug"`
	AllowedHosts []string `yaml:"allowedHosts"`
	globHosts    []glob.Glob
	Database     DbConf
}

func (c *Conf) hostsCompile() {
	if len(c.AllowedHosts) != len(c.globHosts) {
		for _, a := range c.AllowedHosts {
			g, err := glob.Compile(a)
			if err != nil {
				log.WithError(err).WithField("expr", a).Fatal("can't compile allowed host glob")
			}
			c.globHosts = append(c.globHosts, g)
		}
	}
}

// AllowsHost checks peer to determine if it matches a configured host. Returns true if AllowedHosts are not configured.
func (c *Conf) AllowsHost(host string) bool {
	if c.AllowedHosts == nil || len(c.AllowedHosts) == 0 {
		return true
	}
	c.hostsCompile()
	for _, g := range c.globHosts {
		if g.Match(host) {
			return true
		}
	}
	return false
}

// DbDSN is the go-sql connection string for the selected database engine.
// Ex. (for mysql) "user:passwd@tcp(host:port)/database?options"
func (c *Conf) DbDSN() string {
	dsn := ""

	dsn += c.Database.User
	if len(c.Database.Password) > 0 {
		dsn += ":" + c.Database.Password
	}

	dsn += "@"
	if len(c.Database.Socket) > 0 {
		dsn += fmt.Sprintf("unix(%s)", c.Database.Socket)
	} else {
		port := c.Database.Port
		if port == 0 {
			port = 3306
		}
		dsn += fmt.Sprintf("tcp(%s:%d)", c.Database.Host, port)
	}

	if len(c.Database.Database) > 0 {
		dsn += "/" + c.Database.Database
	}

	return dsn
}

// DbDriver is the go-sql database plugin
func (c *Conf) DbDriver() string {
	return c.Database.Driver
}

// GetConf loads and returns the config from file ./config.yaml
func (c *Conf) GetConf(cfgPath string) *Conf {
	yamlFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
