package v1

import (
	"context"
	"exam/api-gateway/api/handlers/models"
	pb "exam/api-gateway/genproto/product-service"
	"exam/api-gateway/pkg/logger"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateProduct
// @Router /v1/product/create [post]
// @Security BearerAuth
// @Summary create product
// @Tags Product
// @Description Insert a new product with provided details
// @Accept json
// @Produce json
// @Param ProductDetails body models.Product true "Create product"
// @Success 201 {object} models.Product
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) CreateProduct(c *gin.Context) {
	var (
		body       models.Product
		jspMarshal protojson.MarshalOptions
	)
	jspMarshal.UseProtoNames = true

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("failed to bind json", logger.Error(err))
		return
	}

	if body.Amount < 0 {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: "amount should be a positive number",
		})
		h.log.Error("user tried to send a negative number", logger.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	resp, err := h.serviceManager.ProductService().CreateProduct(ctx, &pb.Product{
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
		Amount:      body.Amount,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot create product", logger.Error(err))
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Update Product
// @Router /v1/product/update/{id} [put]
// @Security BearerAuth
// @Summary update product
// @Tags Product
// @Description Update ptoduct
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param UserInfo body models.Product true "Update Product"
// @Success 201 {object} models.Product
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) UpdateProduct(c *gin.Context) {
	var (
		body        pb.Product
		jspbMarshal protojson.MarshalOptions
	)
	id := c.Param("id")
	idToInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: err.Error(),
		})
		h.log.Error("cannot convert to int", logger.Error(err))
		return
	}

	jspbMarshal.UseProtoNames = true
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: err.Error(),
		})
		h.log.Error("cannot bind json", logger.Error(err))
		return
	}
	if body.Amount < 0 {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: "amount should be a positive number",
		})
		h.log.Error("user tried to send a negative number", logger.Error(err))
		return
	}

	if body.Price < 0 {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: "price  should be a positive number",
		})
		h.log.Error("user tried to send a negative number", logger.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.ProductService().UpdateProduct(ctx, &pb.Product{
		Id:          int32(idToInt),
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
		Amount:      body.Amount,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot update product", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Get Product By Id
// @Router /v1/product/{id} [get]
// @Security BearerAuth
// @Summary get product by id
// @Tags Product
// @Description Get product
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 201 {object} models.Product
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) GetProductById(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")
	idToInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: err.Error(),
		})
		h.log.Error("cannot convert to int", logger.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.ProductService().GetProductById(ctx, &pb.GetProductId{
		ProductId: int32(idToInt),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot get product", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Delete Product
// @Router /v1/product/delete/{id} [delete]
// @Security BearerAuth
// @Summary delete product
// @Tags Product
// @Description Delete product
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 201 {object} models.Status
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) DeleteProduct(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")
	idToInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: err.Error(),
		})
		h.log.Error("cannot convert to int", logger.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.ProductService().DeleteProduct(ctx, &pb.GetProductId{
		ProductId: int32(idToInt),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})

		h.log.Error("cannot delete product", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Get All Products
// @Router /v1/products/{page}/{limit} [get]
// @Summary get all products
// @Tags Product
// @Description get all products
// @Accept json
// @Produce json
// @Param page path string true "page"
// @Param limit path string true "limit"
// @Success 201 {object} models.ListProducts
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) ListProducts(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	page := c.Param("page")
	pageToInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: err.Error(),
		})
		h.log.Error("cannot parse page query param", logger.Error(err))
		return
	}

	limit := c.Param("limit")
	LimitToInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: err.Error(),
		})
		h.log.Error("cannot parse limit query param", logger.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.ProductService().ListProducts(ctx, &pb.GetListRequest{
		Page:  int32(pageToInt),
		Limit: int32(LimitToInt),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})

		h.log.Error("cannot list products", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Get All Purchased products by user
// @Router /v1/products/get/{id} [get]
// @Security BearerAuth
// @Summary get all purchased products by user id
// @Tags Product
// @Description get all purchased products by user id
// @Accept json
// @Produce json
// @Param page path string true "id"
// @Success 201 {object} models.PurchasedProductsList
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) GetPurchasedProductsByUserId(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	userId := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.ProductService().GetPurchasedProductsByUserId(ctx, &pb.GetUserID{
		UserId: userId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})

		h.log.Error("cannot list products purchased by user", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Buy product
// @Router /v1/product/buy [post]
// @Security BearerAuth
// @Summary buy a product
// @Tags Product
// @Description buy a product
// @Accept json
// @Produce json
// @Param PurchaseInfo body models.BuyProductRequest true "Purchase a product"
// @Success 201 {object} models.BuyProductResponse
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) BuyProduct(c *gin.Context) {
	var (
		res         models.BuyProductResponse
		body        models.BuyProductRequest
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorCodeInvalidJSON,
			Message: err.Error(),
		})
		h.log.Error("cannot bind json", logger.Error(err))
		return
	}

	if body.Amount < 0 {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: "amount should be a positive number",
		})
		h.log.Error("user tried to send a negative number", logger.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	//first check if the product exists
	status, err := h.serviceManager.ProductService().CheckAmount(ctx, &pb.GetProductId{ProductId: body.ProductId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})

		h.log.Error("cannot get list products purchased by user", logger.Error(err))
		return
	}
	if status.Amount == 0 {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: "the product is not currently available, sorry",
		})

		h.log.Error("not available product", logger.Error(err))
		return

	}
	if status.Amount < body.Amount {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: fmt.Sprintf("we have only %d, sorry", status.Amount),
		})

		h.log.Error("not enough", logger.Error(err))
		return
	}
	//Buy a product
	buyResp, err := h.serviceManager.ProductService().BuyProduct(ctx, &pb.BuyProductRequest{
		UserId:    body.UserId,
		ProductId: body.ProductId,
		Amount:    body.Amount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})

		h.log.Error("cannot purchase the product", logger.Error(err))
		return
	}
	//Decrease amount of product from database
	_, err = h.serviceManager.ProductService().DecreaseProductAmount(ctx, &pb.ProductAmountRequest{
		ProductId: body.ProductId,
		AmountBy:  body.Amount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot decrease the amount", logger.Error(err))
		return
	}

	res.Message = "congrats, you've just purchased it!"
	res.ProductId = body.ProductId
	res.UserId = body.UserId
	res.Amount = body.Amount
	res.ProductName = buyResp.Name

	c.JSON(http.StatusOK, res)
}

// Increase the amount of product
// @Router /v1/product/increase [post]
// @Security BearerAuth
// @Summary increase the amount
// @Tags Product
// @Description Increase the amount of product
// @Accept json
// @Produce json
// @Param PurchaseInfo body models.ProductAmountRequest true "Increase the amount"
// @Success 201 {object} models.Status
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) IncreaseAmountOfProduct(c *gin.Context) {
	var (
		body        models.ProductAmountRequest
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorCodeInvalidJSON,
			Message: err.Error(),
		})
		h.log.Error("cannot bind json", logger.Error(err))
		return
	}
	if body.Amount < 0 {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: "amount should be a positive number",
		})
		h.log.Error("user tried to send a negative number", logger.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	_, err = h.serviceManager.ProductService().IncreaseProductAmount(ctx, &pb.ProductAmountRequest{
		ProductId: body.ProductId,
		AmountBy:  body.Amount,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorCodeInvalidJSON,
			Message: err.Error(),
		})
		h.log.Error("cannot increase amount", logger.Error(err))
		return
	}
	c.JSON(http.StatusOK, models.AmountUpdateResp{Success: true, Message: fmt.Sprintf("amount was increased by %d", body.Amount)})
}

// Decrease the amount of product
// @Router /v1/product/decrease [post]
// @Security BearerAuth
// @Summary decrease the amount
// @Tags Product
// @Description Decrease the amount of product
// @Accept json
// @Produce json
// @Param PurchaseInfo body models.ProductAmountRequest true "Decrease the amount"
// @Success 201 {object} models.AmountUpdateResp
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) DecreaseAmountOfProduct(c *gin.Context) {
	var (
		body        models.ProductAmountRequest
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorCodeInvalidJSON,
			Message: err.Error(),
		})
		h.log.Error("cannot bind json", logger.Error(err))
		return
	}

	if body.Amount < 0 {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: "amount should be a positive number",
		})
		h.log.Error("user tried to send a negative number", logger.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	_, err = h.serviceManager.ProductService().DecreaseProductAmount(ctx, &pb.ProductAmountRequest{
		ProductId: body.ProductId,
		AmountBy:  body.Amount,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorCodeInvalidJSON,
			Message: err.Error(),
		})
		h.log.Error("cannot decrease amount", logger.Error(err))
		return
	}
	c.JSON(http.StatusOK, models.AmountUpdateResp{Success: true, Message: fmt.Sprintf("amount was decreased by %d", body.Amount)})
}
