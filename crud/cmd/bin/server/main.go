package main

import (
	"fmt"
	"net/http"
	"time"

	"user-service/config"
	"user-service/handlers"
)

var (
	configPath = "cmd/bin/server/config.yaml"
)

type Config struct {
	Server struct {
		Port    int    `yaml:"port"`
		Addr    string `yaml:"addr"`
		Timeout struct {
			Read  time.Duration `yaml:"read"`
			Write time.Duration `yaml:"write"`
		} `yaml:"timeout"`
	} `yaml:"server"`
}

func main() {
	cfg := &Config{}
	if err := config.Read(configPath, cfg); err != nil {
		panic(fmt.Errorf("failed to read config %q: %v", configPath, err))
	}
	if err := handlers.DumpErrors(); err != nil {
		panic(fmt.Errorf("failed to dump errors: %v", err))
	}
	router, err := handlers.NewRouter()
	if err != nil {
		panic(fmt.Errorf("failed to create router: %v", err))
	}
	server := http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port),
		ReadTimeout:  cfg.Server.Timeout.Read,
		WriteTimeout: cfg.Server.Timeout.Write,
		Handler:      router,
	}
	if err = server.ListenAndServe(); err != nil {
		panic(fmt.Errorf("failed to listen and serve: %v", err))
	}
}
