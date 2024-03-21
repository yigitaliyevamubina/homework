package v1

import (
	t "exam/api-gateway/api/handlers/v1/tokens"
	"exam/api-gateway/config"
	"exam/api-gateway/pkg/logger"
	"exam/api-gateway/services"
	admin "exam/api-gateway/storage/postgresrepo"
	"exam/api-gateway/storage/repo"

	"github.com/casbin/casbin/v2"
)

const (
	ErrorCodeInvalidURL          = "INVALID_URL"
	ErrorCodeInvalidJSON         = "INVALID_JSON"
	ErrorCodeInvalidParams       = "INVALID_PARAMS"
	ErrorCodeInternalServerError = "INTERNAL_SERVER_ERROR"
	ErrorCodeUnauthorized        = "UNAUTHORIZED"
	ErrorCodeAlreadyExists       = "ALREADY_EXISTS"
	ErrorCodeNotFound            = "NOT_FOUND"
	ErrorCodeInvalidCode         = "INVALID_CODE"
	ErrorBadRequest              = "BAD_REQUEST"
	ErrorInvalidCredentials      = "INVALID_CREDENTIALS"
)

type handlerV1 struct {
	inMemoryStorage repo.InMemoryStorageI
	log             logger.Logger
	serviceManager  services.IServiceManager
	cfg             config.Config
	jwtHandler      t.JWTHandler
	postgres        admin.AdminStorageI
	casbin          *casbin.Enforcer
}

type HandlerV1Config struct {
	InMemoryStorage repo.InMemoryStorageI
	Log             logger.Logger
	ServiceManager  services.IServiceManager
	Cfg             config.Config
	JWTHandler      t.JWTHandler
	Postgres        admin.AdminStorageI
	Casbin          *casbin.Enforcer
}

func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		inMemoryStorage: c.InMemoryStorage,
		log:             c.Log,
		serviceManager:  c.ServiceManager,
		cfg:             c.Cfg,
		jwtHandler:      c.JWTHandler,
		postgres:        c.Postgres,
		casbin:          c.Casbin,
	}

}
