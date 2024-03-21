package postgres

import (
	"context"
	"exam/product-service/config"
	pb "exam/product-service/genproto/product-service"
	db2 "exam/product-service/pkg/db"
	"exam/product-service/pkg/logger"
	"exam/product-service/storage/repo"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ProductTestSuite struct {
	suite.Suite
	CleanupFunc func()
	Repository  repo.ProductServiceI
}

func (u *ProductTestSuite) SetupSuite() {
	db, _ := db2.New(*config.Load())
	u.Repository = NewProductRepo(db, logger.New("", ""))
	u.CleanupFunc = db.Close
}

func (u *ProductTestSuite) TestPositionCRUD() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()
	//Create product
	amount := gofakeit.IntRange(5, 50)
	product := &pb.Product{
		Name:        gofakeit.FirstName(),
		Description: gofakeit.ProductDescription(),
		Price:       float32(gofakeit.Price(10.1, 19.2)),
		Amount:      int32(amount),
	}

	createResp, err := u.Repository.CreateProduct(ctx, product)
	u.Suite.NoError(err)
	u.Suite.NotNil(createResp)

	//Get product
	productId := &pb.GetProductId{
		ProductId: createResp.Id,
	}
	getResp, err := u.Repository.GetProductById(ctx, productId)
	u.Suite.NoError(err)
	u.Suite.NotNil(getResp)
	u.Suite.Equal(getResp.Name, createResp.Name)
	u.Suite.Equal(getResp.Amount, createResp.Amount)
	u.Suite.Equal(getResp.Price, createResp.Price)
	u.Suite.Equal(getResp.Description, createResp.Description)

	//List products
	listResp, err := u.Repository.ListProducts(ctx, &pb.GetListRequest{
		Page:  1,
		Limit: 10,
	})
	u.Suite.NoError(err)
	u.Suite.NotNil(listResp)

	//Update product
	updatedName := gofakeit.FirstName()
	product.Name = updatedName
	updatedDescription := gofakeit.ProductDescription()
	product.Description = updatedDescription
	updateResp, err := u.Repository.UpdateProduct(ctx, product)
	u.Suite.NoError(err)
	u.Suite.NotNil(updateResp)
	u.Suite.Equal(updatedName, updateResp.Name)
	u.Suite.Equal(updatedDescription, updateResp.Description)

	//CheckField
	checkResp, err := u.Repository.CheckAmount(ctx, productId)
	u.Suite.NoError(err)
	u.Suite.NotNil(checkResp)
	u.Suite.Equal(checkResp.ProductId, productId.ProductId)

	//Buy product
	userId := uuid.New().String()
	productResp, err := u.Repository.BuyProduct(ctx, &pb.BuyProductRequest{
		UserId:    userId,
		ProductId: productId.ProductId,
		Amount:    1,
	})
	u.Suite.NoError(err)
	u.Suite.NotNil(productResp)
	u.Suite.Equal(productResp.Name, product.Name)
	u.Suite.Equal(productResp.Description, product.Description)

	//Increase product
	response, err := u.Repository.IncreaseProductAmount(ctx, &pb.ProductAmountRequest{
		ProductId: productId.ProductId,
		AmountBy:  1,
	})
	u.Suite.NoError(err)
	u.Suite.NotNil(response)
	u.Suite.Equal(response.IsEnough, true)
	u.Suite.Equal(response.Product.Price, product.Price)

	//Decrease product
	resp, err := u.Repository.DecreaseProductAmount(ctx, &pb.ProductAmountRequest{
		ProductId: productId.ProductId,
		AmountBy:  1,
	})
	u.Suite.NoError(err)
	u.Suite.NotNil(resp)
	u.Suite.Equal(resp.Product.Amount, product.Amount)
	u.Suite.Equal(resp.IsEnough, true)
	u.Suite.Equal(resp.Product.Price, product.Price)

	//Delete product
	_, err = u.Repository.DeleteProduct(ctx, productId)
	u.Suite.NoError(err)
}

func (u *ProductTestSuite) TearDownSuite() {
	u.CleanupFunc()
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}
