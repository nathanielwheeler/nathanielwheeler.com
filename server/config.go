package server

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	configFile = ".config.yml"

	dialect = "postgres"
)

type config struct {
	Env       string         `yaml:"env"`
	Port      int            `yaml:"port"`
	Pepper    string         `yaml:"pepper"`
	HMACKey   string         `yaml:"hmac_key"`
	CSRFBytes int            `yaml:"csrf_bytes"`
	Database  postgresConfig `yaml:"database"`
}

func loadConfig() *config {
	f, err := os.Open(configFile)
	if err != nil {
		panic("No configuration file detected!")
	}
	defer f.Close()
	// Decode file and return Config struct
	var c config
	d := yaml.NewDecoder(f)
	if err := d.Decode(&c); err != nil {
		panic(err)
	}
	return &c
}

func (s *server) isProd() bool {
	return s.config.Env == "prod"
}

// PostgresConfig holds database connection info.
type postgresConfig struct {
	DBName   string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (s *server) dialect() string {
	return dialect
}

// ConnectionString will return a string used to connect to the database
func (s *server) connectionString() string {
	c := s.config.Database
	if s.isProd() {
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s",
			c.Host, c.Port, c.User, c.Password, c.DBName,
		)
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName,
	)
}
