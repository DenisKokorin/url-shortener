package main

import (
	"url-shortener/internal/config"
	"url-shortener/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
}
