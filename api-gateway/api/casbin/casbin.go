package casbin

import (
	"exam/api-gateway/api/handlers/v1/tokens"
	"exam/api-gateway/config"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/cast"
	"net/http"
	"strings"
)

type CasbinHandler struct {
	cfg      config.Config
	enforcer *casbin.Enforcer
}

func NewAuth(casbin *casbin.Enforcer, cfg config.Config) gin.HandlerFunc {
	casbHandler := &CasbinHandler{
		cfg:      cfg,
		enforcer: casbin,
	}

	return func(ctx *gin.Context) {
		allowed, err := casbHandler.CheckPermission(ctx.Request)
		if err != nil {
			v, _ := err.(jwt.ValidationError)
			if v.Errors == jwt.ValidationErrorExpired {
				casbHandler.RequireRefresh(ctx)
				return
			} else {
				casbHandler.RequirePermission(ctx)
				return
			}
		} else if !allowed {
			casbHandler.RequirePermission(ctx)
			return
		}
	}
}

func (a *CasbinHandler) GetRole(ctx *http.Request) (string, int) {
	token := ctx.Header.Get("authorization")
	if token == "" {
		return "unauthorized", http.StatusOK
	}

	var cutToken string

	if strings.Contains(token, "Bearer") {
		cutToken = strings.TrimPrefix(token, "Bearer ")
	} else {
		cutToken = token
	}

	claims, err := tokens.ExtractClaim(cutToken, []byte(a.cfg.SignInKey))
	if err != nil {
		return "unauthorized, token is invalid", http.StatusBadRequest
	}

	return cast.ToString(claims["role"]), http.StatusOK
}

func (a *CasbinHandler) CheckPermission(req *http.Request) (bool, error) {
	role, status := a.GetRole(req)
	if status != http.StatusOK {
		return false, fmt.Errorf(role)
	}

	object := req.URL.Path
	action := req.Method

	response, err := a.enforcer.Enforce(role, object, action)
	if err != nil {
		return false, err
	}

	if !response {
		return false, nil
	}

	return true, nil
}

func (a *CasbinHandler) RequirePermission(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{
		"Status":  "Permission denied",
		"Message": "This method is not allowed to you",
	})
}

func (a *CasbinHandler) RequireRefresh(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"Status":  "unauthorized",
		"Message": "Access token expired",
	})
}
