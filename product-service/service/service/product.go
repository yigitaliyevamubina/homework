package service

import (
	"context"
	pb "exam/product-service/genproto/product-service"
	"exam/product-service/pkg/logger"
	grpcClient "exam/product-service/service/grpc_client"
	"exam/product-service/storage"
)

type ProductService struct {
	storage storage.StorageI
	log     logger.Logger
	service grpcClient.IServiceManager
}

// Constructor
func NewProductService(storage storage.StorageI, log logger.Logger, service grpcClient.IServiceManager) *ProductService {
	return &ProductService{
		storage: storage,
		log:     log,
		service: service,
	}
}
func (c *ProductService) CreateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	return c.storage.ProductService().CreateProduct(ctx, req)
}

func (c *ProductService) GetProductById(ctx context.Context, req *pb.GetProductId) (*pb.Product, error) {
	return c.storage.ProductService().GetProductById(ctx, req)
}

func (c *ProductService) UpdateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	return c.storage.ProductService().UpdateProduct(ctx, req)
}

func (c *ProductService) DeleteProduct(ctx context.Context, req *pb.GetProductId) (*pb.Status, error) {
	return c.storage.ProductService().DeleteProduct(ctx, req)
}

func (c *ProductService) ListProducts(ctx context.Context, req *pb.GetListRequest) (*pb.GetListResponse, error) {
	return c.storage.ProductService().ListProducts(ctx, req)
}

func (c *ProductService) IncreaseProductAmount(ctx context.Context, req *pb.ProductAmountRequest) (*pb.ProductAmountResponse, error) {
	return c.storage.ProductService().IncreaseProductAmount(ctx, req)
}

func (c *ProductService) DecreaseProductAmount(ctx context.Context, req *pb.ProductAmountRequest) (*pb.ProductAmountResponse, error) {
	return c.storage.ProductService().DecreaseProductAmount(ctx, req)
}

func (c *ProductService) CheckAmount(ctx context.Context, req *pb.GetProductId) (*pb.CheckAmountResponse, error) {
	return c.storage.ProductService().CheckAmount(ctx, req)
}

func (c *ProductService) BuyProduct(ctx context.Context, req *pb.BuyProductRequest) (*pb.Product, error) {
	return c.storage.ProductService().BuyProduct(ctx, req)
}

func (c *ProductService) GetPurchasedProductsByUserId(ctx context.Context, req *pb.GetUserID) (*pb.GetPurchasedProductsResponse, error) {
	return c.storage.ProductService().GetPurchasedProductsByUserId(ctx, req)
}
