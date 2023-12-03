package web

import "time"

type ServerConfig struct {
	Port    int    `yaml:"port"`
	Addr    string `yaml:"addr"`
	Timeout struct {
		Read  time.Duration `yaml:"read"`
		Write time.Duration `yaml:"write"`
	} `yaml:"timeout"`
}
