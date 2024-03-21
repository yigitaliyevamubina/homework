package repo

import (
	"context"
	pb "exam/product-service/genproto/product-service"
)

// ProductService interface
type ProductServiceI interface {
	CreateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error)
	GetProductById(ctx context.Context, req *pb.GetProductId) (*pb.Product, error)
	UpdateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error)
	DeleteProduct(ctx context.Context, req *pb.GetProductId) (*pb.Status, error)
	ListProducts(ctx context.Context, req *pb.GetListRequest) (*pb.GetListResponse, error)
	IncreaseProductAmount(ctx context.Context, req *pb.ProductAmountRequest) (*pb.ProductAmountResponse, error)
	DecreaseProductAmount(ctx context.Context, req *pb.ProductAmountRequest) (*pb.ProductAmountResponse, error)
	CheckAmount(ctx context.Context, req *pb.GetProductId) (*pb.CheckAmountResponse, error)
	BuyProduct(ctx context.Context, req *pb.BuyProductRequest) (*pb.Product, error)
	GetPurchasedProductsByUserId(ctx context.Context, req *pb.GetUserID) (*pb.GetPurchasedProductsResponse, error)
}