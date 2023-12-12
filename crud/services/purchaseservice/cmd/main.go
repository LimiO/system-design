package main

import (
	"log"
	"os"

	"onlinestore/internal/config"
	baseweb "onlinestore/pkg/web"
	"onlinestore/services/purchaseservice/web"
)

var (
	configPath = "services/purchaseservice/cmd/config.yaml"
)

type ServiceConfig struct {
	Addr string `yaml:"addr"`
}

type Config struct {
	Server   *baseweb.ServerConfig `yaml:"server"`
	Services struct {
		Payment *ServiceConfig `yaml:"payment"`
		Courier *ServiceConfig `yaml:"courier"`
		Stock   *ServiceConfig `yaml:"stock"`
	} `yaml:"services"`
}

func main() {
	cfg, err := config.Read[Config](configPath)
	if err != nil {
		log.Fatalf("failed to read config %q: %v", configPath, err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("failed to get jwt secret from env JWT_SECRET")
	}
	srv, err := web.NewServer(cfg.Server.Addr, cfg.Server.Port, jwtSecret, cfg.Services.Payment.Addr, cfg.Services.Courier.Addr, cfg.Services.Stock.Addr)
	if err != nil {
		log.Fatalf("failed to make server: %q", err)
	}
	if err = srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
