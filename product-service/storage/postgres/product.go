package postgres

import (
	"context"
	pb "exam/product-service/genproto/product-service"
	"exam/product-service/pkg/db"
	"exam/product-service/pkg/logger"
	"exam/product-service/storage/repo"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
)

type productRepo struct {
	db  *db.Postgres
	log logger.Logger
}

// Constructor
func NewProductRepo(db *db.Postgres, log logger.Logger) repo.ProductServiceI {
	return &productRepo{
		db:  db,
		log: log,
	}
}

func (u *productRepo) CreateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	query := u.db.Builder.Insert("products").
		Columns(`
		name, description, price, amount
		`).
		Values(
			req.Name, req.Description, req.Price, req.Amount,
		).
		Suffix("RETURNING id, created_at")

	err := query.RunWith(u.db.DB).QueryRow().Scan(&req.Id, &req.CreatedAt)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (u *productRepo) GetProductById(ctx context.Context, req *pb.GetProductId) (*pb.Product, error) {
	respProduct := &pb.Product{}

	query := u.db.Builder.Select(`
		id, name, description, price, amount, created_at
	`).From("products").Where(squirrel.Eq{"id": req.ProductId})

	err := query.RunWith(u.db.DB).QueryRow().Scan(
		&respProduct.Id,
		&respProduct.Name,
		&respProduct.Description,
		&respProduct.Price,
		&respProduct.Amount,
		&respProduct.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return respProduct, nil
}

func (u *productRepo) UpdateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	var (
		updateMap = make(map[string]interface{})
		where     = squirrel.And{squirrel.Eq{"id": req.Id}}
	)

	updateMap["name"] = req.Name
	updateMap["description"] = req.Description
	updateMap["price"] = req.Price
	updateMap["amount"] = req.Amount
	updateMap["updated_at"] = time.Now()

	query := u.db.Builder.Update("products").SetMap(updateMap).
		Where(where).
		Suffix("RETURNING updated_at, created_at")

	err := query.RunWith(u.db.DB).QueryRow().Scan(
		&req.UpdatedAt, &req.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (u *productRepo) DeleteProduct(ctx context.Context, req *pb.GetProductId) (*pb.Status, error) {
	query := u.db.Builder.Delete("products").Where(
		squirrel.Eq{"id": req.ProductId},
	)

	_, err := query.RunWith(u.db.DB).Exec()
	if err != nil {
		return &pb.Status{
			Success: false,
		}, err
	}

	return &pb.Status{
		Success: true,
	}, nil
}

func (u *productRepo) ListProducts(ctx context.Context, req *pb.GetListRequest) (*pb.GetListResponse, error) {
	var (
		respProducts = &pb.GetListResponse{Count: 0}
	)

	query := u.db.Builder.Select(
		`id, name, description, price, amount, created_at
	`).From("products")

	query = query.Offset(uint64((req.Page - 1) * req.Limit)).Limit(uint64(req.Limit))

	rows, err := query.RunWith(u.db.DB).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		respProduct := &pb.Product{}
		err = rows.Scan(
			&respProduct.Id,
			&respProduct.Name,
			&respProduct.Description,
			&respProduct.Price,
			&respProduct.Amount,
			&respProduct.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		respProducts.Products = append(respProducts.Products, respProduct)
		respProducts.Count++
	}

	return respProducts, nil
}

func (u *productRepo) IncreaseProductAmount(ctx context.Context, req *pb.ProductAmountRequest) (*pb.ProductAmountResponse, error) {
	var (
		response  = &pb.ProductAmountResponse{Product: &pb.Product{}}
		updateMap = make(map[string]interface{})
		where     = squirrel.And{squirrel.Eq{"id": req.ProductId}}
	)

	var currentAmount int32
	query := u.db.Builder.Select("amount").From("products").Where(squirrel.Eq{"id": req.ProductId})

	err := query.RunWith(u.db.DB).QueryRow().Scan(
		&currentAmount,
	)
	if err != nil {
		return &pb.ProductAmountResponse{
			IsEnough: false,
			Product:  nil,
		}, err
	}

	updateMap["amount"] = req.AmountBy + currentAmount
	updateMap["updated_at"] = time.Now()
	response.Product.Amount = req.AmountBy + currentAmount

	query2 := u.db.Builder.Update("products").SetMap(updateMap).
		Where(where).
		Suffix("RETURNING id, name, description, price, amount, created_at")

	err = query2.RunWith(u.db.DB).QueryRow().Scan(
		&response.Product.Id,
		&response.Product.Name,
		&response.Product.Description,
		&response.Product.Price,
		&response.Product.Amount,
		&response.Product.CreatedAt,
	)
	if err != nil {
		return &pb.ProductAmountResponse{
			IsEnough: false,
			Product:  nil,
		}, err
	}

	response.IsEnough = true

	return response, nil
}

func (u *productRepo) DecreaseProductAmount(ctx context.Context, req *pb.ProductAmountRequest) (*pb.ProductAmountResponse, error) {
	var (
		response  = &pb.ProductAmountResponse{Product: &pb.Product{}}
		updateMap = make(map[string]interface{})
		where     = squirrel.And{squirrel.Eq{"id": req.ProductId}}
	)

	var currentAmount int32
	query := u.db.Builder.Select("amount").From("products").Where(squirrel.Eq{"id": req.ProductId})

	err := query.RunWith(u.db.DB).QueryRow().Scan(
		&currentAmount,
	)
	if err != nil {
		return nil, err
	}

	response.IsEnough = true

	if currentAmount == 0 {
		return nil, fmt.Errorf("not enough")
	}

	// changeTo := currentAmount - req.AmountBy
	// if changeTo < 0 {
	// 	changeTo = currentAmount
	// 	response.IsEnough = false
	// }
	if currentAmount-req.AmountBy < 0 {
		response.IsEnough = false
		return nil, err
	}

	changeTo := currentAmount - req.AmountBy

	updateMap["amount"] = changeTo
	updateMap["updated_at"] = time.Now()

	query2 := u.db.Builder.Update("products").SetMap(updateMap).
		Where(where).
		Suffix("RETURNING id, name, description, price, amount, created_at")

	err = query2.RunWith(u.db.DB).QueryRow().Scan(
		&response.Product.Id,
		&response.Product.Name,
		&response.Product.Description,
		&response.Product.Price,
		&response.Product.Amount,
		&response.Product.CreatedAt,
	)
	if err != nil {
		response.IsEnough = false
		return response, err
	}

	response.Product.Amount = changeTo
	return response, nil
}

func (u *productRepo) CheckAmount(ctx context.Context, req *pb.GetProductId) (*pb.CheckAmountResponse, error) {
	var checkResult pb.CheckAmountResponse
	query := u.db.Builder.Select("amount").From("products").Where(
		squirrel.Eq{"id": req.ProductId},
	)

	err := query.RunWith(u.db.DB).QueryRow().Scan(
		&checkResult.Amount,
	)
	checkResult.ProductId = req.ProductId
	if err != nil {
		return nil, err
	}

	return &checkResult, nil
}

func (u *productRepo) BuyProduct(ctx context.Context, req *pb.BuyProductRequest) (*pb.Product, error) {
	query := u.db.Builder.Insert("users_products").
		Columns("user_id, product_id, amount").
		Values(req.UserId, req.ProductId, req.Amount)

	_, err := query.RunWith(u.db.DB).Exec()
	if err != nil {
		return nil, err
	}

	product, err := u.GetProductById(ctx, &pb.GetProductId{
		ProductId: req.ProductId,
	})

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productRepo) GetPurchasedProductsByUserId(ctx context.Context, req *pb.GetUserID) (*pb.GetPurchasedProductsResponse, error) {
	query := u.db.Builder.Select("product_id").
		From("users_products").
		Where(squirrel.Eq{"user_id": req.UserId})
	rows, err := query.RunWith(u.db.DB).Query()
	if err != nil {
		return nil, err
	}

	var products []*pb.Product

	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		product, err := u.GetProductById(ctx, &pb.GetProductId{ProductId: id})
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

//userId -> many-to-many ([]product_ids) -> for _, id product_id {
//	productInfo, err := GetProductById(id)
//}
