package api

import (
	"github.com/joeshaw/envdecode"
)

var conf *APIConfig

type APIConfig struct {
	Port         string `env:"PORT,default=:8087"`
	DBDriver     string `env:"DB_DRIVER,default=sqlite3"`
	DBConnection string `env:"DB_CONNECTION,default=./data.db"`
}

func LoadConfig() (*APIConfig, error) {
	if conf != nil {
		return conf, nil
	}

	var cfg APIConfig
	if err := envdecode.Decode(&cfg); err != nil {
		return nil, err
	}

	conf = &cfg
	return conf, nil
}
