package api

import (
	"github.com/joeshaw/envdecode"
)

var conf *APIConfig

type APIConfig struct {
	Port       string `env:"PORT,default=:8087"`
	SQLitePath string `env:"SQLITE_PATH,default=./data.db"`
	WithMock   bool   `env:"WITH_MOCK,default=false"`
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
