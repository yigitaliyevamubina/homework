package v1

import (
	"context"
	"encoding/json"
	"exam/api-gateway/api/handlers/models"
	"exam/api-gateway/api/handlers/v1/tokens"
	"exam/api-gateway/email"
	pbp "exam/api-gateway/genproto/product-service"
	pb "exam/api-gateway/genproto/user-service"
	"exam/api-gateway/pkg/utils"

	//pbp "exam/api-gateway/genproto/product-service"
	"exam/api-gateway/pkg/etc"
	"exam/api-gateway/pkg/logger"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// Register User
// @Summary register user
// @Tags User
// @Description Register a new user with provided details
// @Accept json
// @Produce json
// @Param User-data body models.UserRequest true "Register user"
// @Success 201 {object} models.RegisterUserResponse
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
// @Router /v1/user/register [post]
func (h *handlerV1) Register(c *gin.Context) {
	var (
		body       models.UserRequest
		code       string
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

	body.Id = uuid.New().String()
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)
	err = body.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		h.log.Error("error while validating", logger.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	exists, err := h.serviceManager.UserService().CheckField(ctx, &pb.CheckFieldRequest{
		Field: "email",
		Data:  body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("failed to check email uniqueness", logger.Error(err))
		return
	}
	if exists.Status {
		c.JSON(http.StatusConflict, models.ResponseError{
			Code:    ErrorCodeAlreadyExists,
			Message: "This email is already in use. Try another email.",
		})
		h.log.Error("email is not unique", logger.Error(err))
		return
	}
	code = utils.GenerateCode(5)
	respUser := models.RegisterUserModel{
		Id:        body.Id,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Age:       body.Age,
		Email:     body.Email,
		Password:  body.Password,
		Code:      code,
	}

	userJson, err := json.Marshal(respUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInvalidJSON,
			Message: "cannot bind json",
		})
		h.log.Error("cannot bind json", logger.Error(err))
		return
	}

	timeOut := time.Minute * 2

	err = h.inMemoryStorage.SetWithTTL(respUser.Email, string(userJson), int(timeOut.Seconds()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot set with ttl", logger.Error(err))
		return
	}

	from := "mubinayigitaliyeva00@gmail.com"
	password := "iocd vnhb lnvx digm"
	err = email.SendVerificationCode(email.Params{
		From:     from,
		Password: password,
		To:       respUser.Email,
		Message:  fmt.Sprintf("Hi %s,", respUser.FirstName),
		Code:     respUser.Code,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot send a code to an email", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.RegisterUserResponse{Message: "a verification code was sent to your email, please check it."})
}

// Verify User
// @Summary verify user
// @Tags User
// @Description Verify a user with code sent to their email
// @Accept json
// @Product json
// @Param email path string true "email"
// @Param code path string true "code"
// @Success 201 {object} models.VerifyUserResponse
// @Failure 400 string error models.ResponseError
// @Failure 400 string error models.ResponseError
// @Router /v1/user/verify/{email}/{code} [get]
func (h *handlerV1) Verify(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	userEmail := c.Param("email")
	code := c.Param("code")

	user, err := redis.Bytes(h.inMemoryStorage.Get(userEmail))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorCodeUnauthorized,
			Message: "code is expired, try again.",
		})
		h.log.Error("Code is expired, TTL is over.")
		return
	}

	var respUser models.RegisterUserModel
	if err := json.Unmarshal(user, &respUser); err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot unmarshal user from redis", logger.Error(err))
		fmt.Println(respUser)
		return
	}

	if respUser.Code != code {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorCodeInvalidCode,
			Message: "code is incorrect, try again.",
		})
		h.log.Error("verification failed", logger.Error(err))
		return
	}

	respUser.Password, err = etc.HashPassword(respUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot hash the password", logger.Error(err))
		return
	}

	h.jwtHandler = tokens.JWTHandler{
		Sub:       respUser.Id,
		Role:      "user",
		SignInKey: h.cfg.SignInKey,
		Log:       h.log,
		Timeout:   h.cfg.AccessTokenTimeout,
	}

	access, refresh, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot create access and refresh token", logger.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	_, err = h.serviceManager.UserService().CreateUser(ctx, &pb.User{
		Id:           respUser.Id,
		FirstName:    respUser.FirstName,
		LastName:     respUser.LastName,
		Age:          respUser.Age,
		Email:        respUser.Email,
		Password:     respUser.Password,
		RefreshToken: refresh,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot create user", logger.Error(err))
		return
	}

	res := models.UserModel{
		Id:          respUser.Id,
		FirstName:   respUser.FirstName,
		LastName:    respUser.LastName,
		Age:         respUser.Age,
		Email:       respUser.Email,
		Password:    respUser.Password,
		AccessToken: access,
	}

	c.JSON(http.StatusOK, res)
}

// Login User
// @Summary login user
// @Tags User
// @Description Login
// @Accept json
// @Produce json
// @Param User body models.LoginRequest true "Login"
// @Success 201 {object} models.UserModel
// @Failure 400 string Error models.ResponseError
// @Failure 400 string Error models.ResponseError
// @Router /v1/user/login [post]
func (h *handlerV1) Login(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
		body        models.LoginRequest
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	resp, err := h.serviceManager.UserService().Check(ctx, &pb.IfExists{
		Email: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot get user by email", logger.Error(err))
		return
	}

	if !etc.CompareHashPassword(resp.Password, body.Password) {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: "invalid code, try again",
		})
		h.log.Error("wrong password", logger.Error(err))
		return
	}
	h.jwtHandler = tokens.JWTHandler{
		Sub:       resp.Id,
		Role:      "user",
		SignInKey: h.cfg.SignInKey,
		Log:       h.log,
		Timeout:   h.cfg.AccessTokenTimeout,
	}

	access, _, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot create access token", logger.Error(err))
		return
	}

	res := models.UserModel{
		Id:          resp.Id,
		FirstName:   resp.FirstName,
		LastName:    resp.LastName,
		Age:         resp.Age,
		Email:       resp.Email,
		Password:    resp.Password,
		AccessToken: access,
	}

	c.JSON(http.StatusOK, res)
}

// CreateUser
// @Router /v1/user/create [post]
// @Security BearerAuth
// @Summary create user
// @Tags User
// @Description Create a new user with provided details
// @Accept json
// @Produce json
// @Param UserInfo body models.UserRequest true "Create user"
// @Success 201 {object} models.UserModel
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		body        models.UserRequest
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

	body.Password, err = etc.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot hash the password", logger.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	body.Id = uuid.New().String()
	h.jwtHandler = tokens.JWTHandler{
		Sub:       body.Id,
		Role:      "user",
		SignInKey: h.cfg.SignInKey,
		Log:       h.log,
		Timeout:   h.cfg.AccessTokenTimeout,
	}

	access, refresh, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot generate refresh token", logger.Error(err))
		return
	}

	response, err := h.serviceManager.UserService().CreateUser(ctx, &pb.User{
		Id:           body.Id,
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		Age:          body.Age,
		Email:        body.Email,
		Password:     body.Password,
		RefreshToken: refresh,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot create user", logger.Error(err))
		return
	}

	res := models.UserModel{
		Id:          response.Id,
		FirstName:   response.FirstName,
		LastName:    response.LastName,
		Age:         response.Age,
		Email:       response.Email,
		Password:    response.Email,
		AccessToken: access,
	}
	c.JSON(http.StatusCreated, res)
}

// Get User By Id
// @Router /v1/user/{id} [get]
// @Security BearerAuth
// @Summary get user by id
// @Tags User
// @Description Get user
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 201 {object} models.UserWithProducts
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) GetUserById(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()
	response, err := h.serviceManager.UserService().GetUserById(ctx, &pb.GetUserId{
		UserId: id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot get user", logger.Error(err))
		return
	}

	resp := models.UserWithProducts{
		User: models.UserRequest{
			Id:        response.Id,
			FirstName: response.FirstName,
			LastName:  response.LastName,
			Age:       response.Age,
			Email:     response.Email,
			Password:  response.Password,
		},
		Products: []*models.Product{},
	}

	prods, err := h.serviceManager.ProductService().GetPurchasedProductsByUserId(ctx, &pbp.GetUserID{UserId: response.Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot get purchased products of a user", logger.Error(err))
		return
	}
	for _, prod := range prods.Products {
		resp.Products = append(resp.Products, &models.Product{
			Id:          prod.Id,
			Name:        prod.Name,
			Description: prod.Description,
			Price:       prod.Price,
			Amount:      prod.Amount,
			Created_at:  prod.CreatedAt,
			Updated_at:  prod.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// Update User
// @Router /v1/user/update/{id} [put]
// @Security BearerAuth
// @Summary update user
// @Tags User
// @Description Update user
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param UserInfo body models.VerifyUserResponse true "Update User"
// @Success 201 {object} models.VerifyUserResponse
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var (
		body        pb.User
		jspbMarshal protojson.MarshalOptions
	)
	id := c.Param("id")

	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: err.Error(),
		})
		h.log.Error("cannot bind json", logger.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.UserService().UpdateUser(ctx, &pb.User{
		Id:        id,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Age:       body.Age,
		Email:     body.Email,
		Password:  body.Password,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})
		h.log.Error("cannot update user", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Delete User
// @Router /v1/user/delete/{id} [delete]
// @Security BearerAuth
// @Summary delete user
// @Tags User
// @Description Delete user
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 201 {object} models.Status
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) DeleteUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.UserService().DeleteUser(ctx, &pb.GetUserId{
		UserId: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:    ErrorCodeInternalServerError,
			Message: err.Error(),
		})

		h.log.Error("cannot delete user", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Get All Users
// @Router /v1/users/{page}/{limit} [get]
// @Security BearerAuth
// @Summary get all users
// @Tags User
// @Description get all users
// @Accept json
// @Produce json
// @Param page path string true "page"
// @Param limit path string true "limit"
// @Success 201 {object} models.ListUsers
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) GetAllUsers(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()
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

	response, err := h.serviceManager.UserService().ListUsers(ctx, &pb.GetListRequest{
		Page:  int32(pageToInt),
		Limit: int32(LimitToInt),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot get all users", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Update Access token
// @Router /v1/user/refresh [post]
// @Security BearerAuth
// @Summary update access token
// @Tags User
// @Description get access token updated
// @Accept json
// @Produce json
// @Param refresh-token body models.AccessTokenUpdateReq true "refresh-token"
// @Success 201 {object} models.AccessTokenUpdateResp
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
func (h *handlerV1) UpdateRefreshToken(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
		body        models.AccessTokenUpdateReq
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:    ErrorBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	claims, err := tokens.ExtractClaim(body.RefreshToken, []byte(h.cfg.SignInKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	userId := cast.ToString(claims["sub"])

	h.jwtHandler = tokens.JWTHandler{
		Sub:       userId,
		Role:      "user",
		SignInKey: h.cfg.SignInKey,
		Log:       h.log,
		Timeout:   h.cfg.AccessTokenTimeout,
	}

	access, refresh, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = h.serviceManager.UserService().UpdateRefreshToken(ctx, &pb.UpdateRefreshTokenReq{
		UserId:       cast.ToString(claims["sub"]),
		RefreshToken: refresh,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp := models.AccessTokenUpdateResp{
		Status:      "success",
		UserID:      userId,
		AccessToken: access,
	}

	c.JSON(http.StatusOK, resp)
}
