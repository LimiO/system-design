package main

import (
	"log"
	"os"

	"onlinestore/internal/config"
	baseweb "onlinestore/pkg/web"
	"onlinestore/services/backend/web"
)

var (
	configPath = "services/backend/cmd/config.yaml"
)

type ServiceConfig struct {
	Addr string `yaml:"addr"`
}

type Config struct {
	Server   *baseweb.ServerConfig `yaml:"server"`
	Services struct {
		Payment       *ServiceConfig `yaml:"payment"`
		Authorization *ServiceConfig `yaml:"authorization"`
		Purchases     *ServiceConfig `yaml:"purchases"`
		User          *ServiceConfig `yaml:"user"`
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
	srv, err := web.NewServer(
		cfg.Server.Addr, cfg.Server.Port, jwtSecret,
		cfg.Services.Payment.Addr,
		cfg.Services.Authorization.Addr,
		cfg.Services.Purchases.Addr,
		cfg.Services.User.Addr,
	)
	if err != nil {
		log.Fatalf("failed to make server: %q", err)
	}
	log.Println("start server!")
	if err = srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
