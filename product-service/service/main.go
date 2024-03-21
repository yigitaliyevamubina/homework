package service

import (
	"context"
	"exam/product-service/config"
	pb "exam/product-service/genproto/product-service"
	// "exam/product-service/pkg/db"
	"exam/product-service/pkg/logger"
	grpcClient2 "exam/product-service/service/grpc_client"
	"exam/product-service/service/service"
	storage2 "exam/product-service/storage"
	"fmt"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type Service struct {
	ProductService *service.ProductService
}

func New(cfg *config.Config, log logger.Logger) (*Service, error) {
	// postgres, err := db.New(*cfg)
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot connect to database:", err.Error())
	// }
	// storage := storage2.New(postgres, log)

	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	database := client.Database("productdb")
	storage := storage2.New(database, log)
	grpcClient, err := grpcClient2.New(*cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to grpc client:%v", err.Error())
	}

	return &Service{ProductService: service.NewProductService(storage, log, grpcClient)}, nil
}

func (s *Service) Run(log logger.Logger, cfg *config.Config) {
	server := grpc.NewServer()

	pb.RegisterProductServiceServer(server, s.ProductService)

	listen, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("error while creating a listener", logger.Error(err))
		return
	}

	defer logger.Cleanup(log)

	log.Info("main: sqlConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase),
		logger.String("rpc port", cfg.RPCPort))

	if err := server.Serve(listen); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
