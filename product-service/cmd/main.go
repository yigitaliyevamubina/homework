package main

import (
	"exam/product-service/config"
	"exam/product-service/pkg/logger"
	service2 "exam/product-service/service"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.Environment, "product-service")
	service, err := service2.New(cfg, log)
	if err != nil {
		log.Error("error while accessing services", logger.Error(err))
		return
	}

	service.Run(log, cfg)
}
