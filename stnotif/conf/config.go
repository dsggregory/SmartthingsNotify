package conf

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Conf an instantiated configuration
type Conf struct {
	ServerPort int      `yaml:"serverPort"`
	Debug      bool     `yaml:"debug"`
	Hosts      []string `yaml:"hosts"`
	// DbDriver is the go-sql database plugin
	DbDriver string `yaml:"dbDriver"`
	// DbDSN is the go-sql connection string for the selected database engine.
	// Ex. (for mysql) "user:passwd@tcp(host:port)/database?options"
	DbDSN string `yaml:"dbDSN"`
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
