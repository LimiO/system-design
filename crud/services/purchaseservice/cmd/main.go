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

type Config struct {
	Server *baseweb.ServerConfig `yaml:"server"`
}

func main() {
	// TODO запустить горутину, которая через 15 минут отменит заказ
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
	log.Println("start server!")
	if err = srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
