package grpcClient

import (
	"exam/user-service/config"
	pbp "exam/user-service/genproto/product-service"
	"exam/user-service/pkg/logger"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type IServiceManager interface {
	ProductService() pbp.ProductServiceClient
}

type serviceManager struct {
	cfg            config.Config
	productService pbp.ProductServiceClient
}

func New(cfg config.Config) (IServiceManager, error) {
	connProduct, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.ProductServiceHost, cfg.ProductServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error while dialing to the product service", logger.Error(err))
	}
	return &serviceManager{
		cfg:            cfg,
		productService: pbp.NewProductServiceClient(connProduct),
	}, nil
}

func (s *serviceManager) ProductService() pbp.ProductServiceClient {
	return s.productService
}
