package storage

import (
	// "exam/product-service/pkg/db"
	"exam/product-service/pkg/logger"
	// "exam/product-service/storage/postgres"
	mon "exam/product-service/storage/mongo"
	"exam/product-service/storage/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

// Storage
type StorageI interface {
	ProductService() repo.ProductServiceI
}

type storagePg struct {
	productService repo.ProductServiceI
}

func New(db *mongo.Database, log logger.Logger) StorageI {
	return &storagePg{productService: mon.NewProductRepo(db, log)}
}

func (s *storagePg) ProductService() repo.ProductServiceI {
	return s.productService
}
