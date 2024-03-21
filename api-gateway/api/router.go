package api

import (
	// casb "exam/api-gateway/api/casbin"
	_ "exam/api-gateway/api/docs"
	v1 "exam/api-gateway/api/handlers/v1"
	"exam/api-gateway/api/handlers/v1/tokens"
	"exam/api-gateway/config"
	"exam/api-gateway/pkg/logger"
	"exam/api-gateway/services"
	"exam/api-gateway/storage/repo"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"fmt"
	gormadapter "github.com/casbin/gorm-adapter/v3"

	admin "exam/api-gateway/storage/postgresrepo"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	"github.com/gin-gonic/gin"
	//swaggerFiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
)

// Option Struct
type Option struct {
	InMemory       repo.InMemoryStorageI
	Cfg            config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Postgres       admin.AdminStorageI
}

// New -> constructor
// @title THIRD EXAM
// @version 1.0
// @description Auth, Role-management, Product, User
// @host localhost:7070
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func New(option Option) *gin.Engine {
	psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		option.Cfg.PostgresHost,
		option.Cfg.PostgresPort,
		option.Cfg.PostgresUser,
		option.Cfg.PostgresPassword,
		option.Cfg.PostgresDatabase)

	adapter, err := gormadapter.NewAdapter("postgres", psqlString, true)
	if err != nil {
		option.Logger.Fatal("error while updating new adapter", logger.Error(err))
	}

	casbinEnforcer, err := casbin.NewEnforcer(option.Cfg.AuthConfigPath, adapter)
	if err != nil {
		option.Logger.Error("cannot create a new enforcer", logger.Error(err))
	}

	err = casbinEnforcer.LoadPolicy()
	if err != nil {
		panic(err)
	}

	casbinEnforcer.GetRoleManager().AddMatchingFunc("keyMatch", util.KeyMatch)
	casbinEnforcer.GetRoleManager().AddMatchingFunc("keyMatch3", util.KeyMatch3)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	jwtHandle := tokens.JWTHandler{
		SignInKey: option.Cfg.SignInKey,
		Log:       option.Logger,
	}

	handlerV1 := v1.New(&v1.HandlerV1Config{
		InMemoryStorage: option.InMemory,
		Log:             option.Logger,
		ServiceManager:  option.ServiceManager,
		Cfg:             option.Cfg,
		JWTHandler:      jwtHandle,
		Postgres:        option.Postgres,
		Casbin:          casbinEnforcer,
	})

	api := router.Group("/v1")

	// api.Use(casb.NewAuth(casbinEnforcer, option.Cfg))

	//rbac
	api.GET("/rbac/roles", handlerV1.ListAllRoles)              //superadmin
	api.GET("/rbac/policies/:role", handlerV1.ListRolePolicies) //superadmin
	api.POST("/rbac/add/policy", handlerV1.AddPolicyToRole)     //superadmin
	api.DELETE("/rbac/delete/policy", handlerV1.DeletePolicy)   //superadmin

	//users
	api.POST("/user/create", handlerV1.CreateUser)         //admin
	api.POST("/user/register", handlerV1.Register)         //unauthorized
	api.GET("/user/:id", handlerV1.GetUserById)            //user
	api.PUT("/user/update/:id", handlerV1.UpdateUser)      //user
	api.DELETE("/user/delete/:id", handlerV1.DeleteUser)   //admin
	api.GET("/users/:page/:limit", handlerV1.GetAllUsers)  //admin
	api.GET("/user/verify/:email/:code", handlerV1.Verify) //unauthorized
	api.POST("/user/login", handlerV1.Login)               //unauthorized
	api.POST("/user/refresh", handlerV1.UpdateRefreshToken)

	//product
	api.POST("/product/create", handlerV1.CreateProduct)                 //admin
	api.PUT("/product/update/:id", handlerV1.UpdateProduct)              //admin
	api.GET("/product/:id", handlerV1.GetProductById)                    //user
	api.DELETE("/product/delete/:id", handlerV1.DeleteProduct)           //admin
	api.GET("/products/:page/:limit", handlerV1.ListProducts)            //unauthorized
	api.GET("/products/get/:id", handlerV1.GetPurchasedProductsByUserId) //user
	api.POST("/product/buy", handlerV1.BuyProduct)                       //user
	api.POST("/product/increase", handlerV1.IncreaseAmountOfProduct)     //admin
	api.POST("/product/decrease", handlerV1.DecreaseAmountOfProduct)     //admin

	//admin
	api.POST("/auth/create", handlerV1.CreateAdmin)   //superadmin
	api.DELETE("/auth/delete", handlerV1.DeleteAdmin) //superadmin
	api.POST("/auth/login", handlerV1.LoginAdmin)     //admin

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return router
}
