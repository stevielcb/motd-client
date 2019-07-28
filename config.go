package main

import (
	"github.com/kelseyhightower/envconfig"
)

var (
	c Config
)

type Config struct {
	Host      string `default:"localhost"`
	Port      int    `default:"4200"`
	TimeoutMs int    `default:"100" split_words:"true"`
}

func init() {
	err := envconfig.Process("motd", &c)
	if err != nil {
		panic(err)
	}
}
