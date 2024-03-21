package main

import (
	"exam/user-service/config"
	"exam/user-service/pkg/logger"
	service2 "exam/user-service/service"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.Environment, "user-service")
	service, err := service2.New(cfg, log)
	if err != nil {
		log.Error("error while accessing services", logger.Error(err))
		return
	}

	service.Run(log, cfg)
}
