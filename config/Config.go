package config

import (
	"flag"
	"github.com/naoina/toml"
	"os"
)

type Config struct {
	ApiKey string
}

func SetUp() *Config {
	var pathFlag = flag.String("config", "./key.toml", "configuration setting")
	flag.Parse()
	cfg := new(Config)
	file, _ := os.Open(*pathFlag)
	if err := toml.NewDecoder(file).Decode(cfg); err != nil {
		panic(err)
	}
	return cfg
}
