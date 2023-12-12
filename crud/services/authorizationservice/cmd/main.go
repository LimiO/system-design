package main

import (
	"log"
	"os"

	"onlinestore/internal/config"
	baseweb "onlinestore/pkg/web"
	"onlinestore/services/authorizationservice/web"
)

var (
	configPath = "services/authorizationservice/cmd/config.yaml"
)

type Config struct {
	Server *baseweb.ServerConfig `yaml:"server"`
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
	srv, err := web.NewServer(cfg.Server.Addr, cfg.Server.Port, jwtSecret)
	if err != nil {
		log.Fatalf("failed to make server: %q", err)
	}
	if err = srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
