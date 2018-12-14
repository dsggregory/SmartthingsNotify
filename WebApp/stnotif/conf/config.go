package conf

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

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
	ServerPort int      `yaml:"serverPort"`
	Debug      bool     `yaml:"debug"`
	Hosts      []string `yaml:"hosts"`
	Database   DbConf
}

// AllowsHost checks host against to determine if a peer matches a configured host.
func (c *Conf) AllowsHost(host string) bool {
	for _, a := range c.Hosts {
		if a == host {
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
