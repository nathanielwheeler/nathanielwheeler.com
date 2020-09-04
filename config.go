package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	devConfigFile  = ".dev_config.yml"
	prodConfigFile = ".prod_config.yml"

	dialect = "postgres"
)

// Config holds configuration variables
type Config struct {
	Env       string         `yaml:"env"`
	Port      int            `yaml:"port"`
	Pepper    string         `yaml:"pepper"`
	HMACKey   string         `yaml:"hmac_key"`
	CSRFBytes int            `yaml:"csrf_bytes"`
	Database  PostgresConfig `yaml:"database"`
}

// LoadConfig will load production or development configuration files.
func LoadConfig() Config {
	f, err := os.Open(prodConfigFile)
	if err != nil {
		f, err = os.Open(devConfigFile)
		if err != nil {
			panic("No configuration file detected!")
		}
		fmt.Println("Using DEVELOPMENT configuration.")
	} else {
		fmt.Println("Using PRODUCTION configuration.")
  }
  defer f.Close()
	// Decode file and return Config struct
  var c Config
  d := yaml.NewDecoder(f)
	if err := d.Decode(&c); err != nil {
    panic(err)
  }
	return c
}

// IsProd sets Config Env to Production
func (c Config) IsProd() bool {
	return c.Env == "prod"
}

// PostgresConfig holds database connection info.
type PostgresConfig struct {
	DBName   string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Dialect will return the dialect that GORM will use.
func (c PostgresConfig) Dialect() string {
	return dialect
}

// ConnectionString will return a string used to connect to the database
func (c PostgresConfig) ConnectionString() string {
	if c.Password == "" {
		return fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.DBName,
		)
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName,
	)
}
