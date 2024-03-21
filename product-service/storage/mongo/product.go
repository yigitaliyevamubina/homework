package mongo

import (
	"context"
	pb "exam/product-service/genproto/product-service"
	"exam/product-service/pkg/logger"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productRepo struct {
	database *mongo.Database
	log      logger.Logger
}

func NewProductRepo(database *mongo.Database, log logger.Logger) *productRepo {
	return &productRepo{database: database, log: log}
}

func (p *productRepo) CreateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	collection := p.database.Collection("products")
	result, err := collection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}

	var response pb.Product
	filter := bson.M{"_id": result.InsertedID}
	err = collection.FindOne(ctx, filter).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (p *productRepo) GetProductById(ctx context.Context, req *pb.GetProductId) (*pb.Product, error) {
	collection := p.database.Collection("products")

	var response pb.Product
	filter := bson.M{"id": req.ProductId}
	err := collection.FindOne(ctx, filter).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (p *productRepo) UpdateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	collection := p.database.Collection("products")

	var response pb.Product

	filter := bson.M{"id": req.Id}

	updateReq := bson.M{
		"$set": bson.M{
			"name":        req.Name,
			"description": req.Description,
			"price":       req.Price,
			"amount":      req.Amount,
			"updated_at":  time.Now(),
		},
	}

	err := collection.FindOneAndUpdate(ctx, filter, updateReq).Decode(&req)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (p *productRepo) DeleteProduct(ctx context.Context, req *pb.GetProductId) (*pb.Status, error) {
	collection := p.database.Collection("products")

	filter := bson.M{"id": req.ProductId}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return &pb.Status{Success: false}, err
	}

	return &pb.Status{Success: true}, nil
}

func (p *productRepo) ListProducts(ctx context.Context, req *pb.GetListRequest) (*pb.GetListResponse, error) {
	collection := p.database.Collection("products")

	var response pb.GetListResponse

	reqOptions := options.Find()

	reqOptions.SetSkip(int64(req.Page-1) * int64(req.Limit))
	reqOptions.SetLimit(int64(req.Limit))

	cursor, err := collection.Find(ctx, bson.M{}, reqOptions)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var product pb.Product
		err = cursor.Decode(&product)
		if err != nil {
			return nil, err
		}

		response.Count++
		response.Products = append(response.Products, &product)
	}

	return &response, nil
}

func (p *productRepo) IncreaseProductAmount(ctx context.Context, req *pb.ProductAmountRequest) (*pb.ProductAmountResponse, error) {
	collection := p.database.Collection("products")

	var response pb.Product
	filter := bson.M{"id": req.ProductId}
	err := collection.FindOne(ctx, filter).Decode(&response)
	if err != nil {
		return &pb.ProductAmountResponse{IsEnough: false, Product: nil}, err
	}

	updateReq := bson.M{
		"$set": bson.M{
			"amount":     req.AmountBy + response.Amount,
			"updated_at": time.Now(),
		},
	}

	err = collection.FindOneAndUpdate(ctx, filter, updateReq).Decode(&req)
	if err != nil {
		return nil, err
	}

	return &pb.ProductAmountResponse{IsEnough: true, Product: &response}, nil
}

func (p *productRepo) DecreaseProductAmount(ctx context.Context, req *pb.ProductAmountRequest) (*pb.ProductAmountResponse, error) {
	collection := p.database.Collection("products")

	var response pb.Product
	filter := bson.M{"id": req.ProductId}
	err := collection.FindOne(ctx, filter).Decode(&response)
	if err != nil {
		return &pb.ProductAmountResponse{IsEnough: false, Product: nil}, err
	}

	if response.Amount == 0 {
		return nil, fmt.Errorf("not enough")
	}

	if response.Amount-req.AmountBy < 0 {
		return &pb.ProductAmountResponse{IsEnough: false, Product: &response}, err
	}

	updateReq := bson.M{
		"$set": bson.M{
			"amount":     response.Amount - req.AmountBy,
			"updated_at": time.Now(),
		},
	}

	err = collection.FindOneAndUpdate(ctx, filter, updateReq).Decode(&req)
	if err != nil {
		return nil, err
	}

	return &pb.ProductAmountResponse{IsEnough: true, Product: &response}, nil
}

func (p *productRepo) CheckAmount(ctx context.Context, req *pb.GetProductId) (*pb.CheckAmountResponse, error) {
	collection := p.database.Collection("products")

	var checkResult pb.CheckAmountResponse
	var response pb.Product
	filter := bson.M{"id": req.ProductId}
	err := collection.FindOne(ctx, filter).Decode(&response)
	if err != nil {
		return nil, err
	}

	checkResult.Amount = response.Amount
	checkResult.ProductId = response.Id

	return &checkResult, nil
}

func (p *productRepo) BuyProduct(ctx context.Context, req *pb.BuyProductRequest) (*pb.Product, error) {
	collection := p.database.Collection("users_products")

	_, err := collection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}

	response, err := p.GetProductById(ctx, &pb.GetProductId{ProductId: req.ProductId})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *productRepo) GetPurchasedProductsByUserId(ctx context.Context, req *pb.GetUserID) (*pb.GetPurchasedProductsResponse, error) {
	collection := p.database.Collection("users_products")

	var products []*pb.Product

	cursor, err := collection.Find(ctx, bson.M{"user_id": req.UserId})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var order pb.BuyProductRequest
		err := cursor.Decode(&order)
		if err != nil {
			return nil, err
		}

		product, err := p.GetProductById(ctx, &pb.GetProductId{ProductId: order.ProductId})
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	response := &pb.GetPurchasedProductsResponse{
		Products: products,
	}

	return response, nil
}
