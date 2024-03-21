package main

import (
	"exam/api-gateway/api"
	"exam/api-gateway/config"
	"exam/api-gateway/pkg/db"
	"exam/api-gateway/pkg/etc"
	"exam/api-gateway/pkg/logger"
	"exam/api-gateway/services"
	admin "exam/api-gateway/storage/postgres"
	"exam/api-gateway/storage/redis"
	"fmt"

	rds "github.com/gomodule/redigo/redis"
)

func main() {
	//login superadmin -> username = 'a' -> password = 'b'
	fmt.Println(etc.HashPassword("admin"))
				
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api_gateway")

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error", logger.Error(err))
	}
	fmt.Println(etc.HashPassword("b"))


	redisPool := rds.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (rds.Conn, error) {
			c, err := rds.Dial("tcp", fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	db, _, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("cannot run http server", logger.Error(err))
		panic(err)
	}

	server := api.New(api.Option{
		InMemory:       redis.NewRedisRepo(&redisPool),
		Cfg:            cfg,
		Logger:         log,
		ServiceManager: serviceManager,
		Postgres:       admin.NewAdminRepo(db),
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("cannot run http server", logger.Error(err))
		panic(err)
	}
}
