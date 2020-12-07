package server

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
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
	f, err := os.Open(".config.yml")
	if err != nil {
		// This error can be hit during testing.  If so, it means I didn't set the working directory in the test.
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

type postgresConfig struct {
	DBName   string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Returns a connection string representing a URI in the form of:
// postgresql://[user[:password]@][port][:port][/dbname]
func (s *server) connectionString() string {
	c := s.config.Database

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}
