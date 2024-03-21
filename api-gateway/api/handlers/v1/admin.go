package v1

import (
	"exam/api-gateway/api/handlers/models"
	"exam/api-gateway/api/handlers/v1/tokens"
	"exam/api-gateway/pkg/etc"
	"exam/api-gateway/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

// Create Admin
// @Router /v1/auth/create [post]
// @Security BearerAuth
// @Summary create admin
// @Tags Auth
// @Description Create a new admin if you are a superadmin
// @Accept json
// @Product json
// @Param super-username query string true "super-username"
// @Param super-password query string true "super-password"
// @Param admin body models.AdminReq true "admin"
// @Success 201 {object} models.SuperAdminMessage
// @Failure 400 string error models.Error
// @Failure 400 string error models.Error
func (h *handlerV1) CreateAdmin(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
		body        models.AdminReq
	)
	jspbMarshal.UseProtoNames = true

	superAdminUsername := c.Query("super-username")
	superAdminPassword := c.Query("super-password")

	if superAdminUsername == "a" && superAdminPassword == "b" {
		err := c.BindJSON(&body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to bind json", logger.Error(err))
			return
		}

		body.Id = uuid.NewString()

		body.Password, err = etc.HashPassword(body.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("cannot hash the password", logger.Error(err))
			return
		}

		h.jwtHandler = tokens.JWTHandler{
			Sub:       body.Id,
			Role:      "admin",
			SignInKey: h.cfg.SignInKey,
			Log:       h.log,
			Timeout:   h.cfg.AccessTokenTimeout,
		}

		_, refresh, err := h.jwtHandler.GenerateAuthJWT()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			h.log.Error("cannot create new access and refresh token", logger.Error(err))
			return
		}

		adminResp := models.AdminResp{
			Id:           body.Id,
			FullName:     body.FullName,
			Age:          body.Age,
			UserName:     body.UserName,
			Email:        body.Email,
			Password:     body.Password,
			RefreshToken: refresh,
		}

		err = h.postgres.Create(&adminResp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("cannot create admin", logger.Error(err))
			return
		}

		c.JSON(http.StatusCreated, models.SuperAdminMessage{
			Message: "admin successfully created",
		})
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "you cannot create admin, provide username and password",
		})
	}
}

// Delete Admin
// @Router /v1/auth/delete [delete]
// @Security BearerAuth
// @Summary delete admin
// @Tags Auth
// @Description delete admin if you are a superadmin
// @Accept json
// @Product json
// @Param super-username query string true "super-username"
// @Param super-password query string true "super-password"
// @Param admin body models.DeleteAdmin true "admin"
// @Success 201 {object} models.SuperAdminMessage
// @Failure 400 string error models.Error
// @Failure 400 string error models.Error
func (h *handlerV1) DeleteAdmin(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
		body        models.DeleteAdmin
	)
	jspbMarshal.UseProtoNames = true

	superAdminUsername := c.Query("super-username")
	superAdminPassword := c.Query("super-password")

	if superAdminUsername == "a" && superAdminPassword == "b" {
		err := c.BindJSON(&body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to bind json", logger.Error(err))
			return
		}

		_, password, status, err := h.postgres.Get(body.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("cannot get admin", logger.Error(err))
			return
		}

		if !status {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "this admin does not exist",
			})
			h.log.Error("admin does not exist", logger.Error(err))
			return
		}

		if !etc.CompareHashPassword(password, body.Password) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "incorrect password",
			})
			h.log.Error("incorrect password", logger.Error(err))
			return
		}

		resp := h.postgres.Delete(body.Username, body.Password)
		if resp != nil {
			if resp.Error() == "no rows were deleted" {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "this admin does not exists",
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": resp.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.SuperAdminMessage{
			Message: "admin successfully deleted",
		})
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "permission denied",
			"message": "provide correct username and password for superadmin",
		})
		return
	}

}

// Delete Admin
// @Summary login
// @Tags Auth
// @Description login as admin
// @Accept json
// @Product json
// @Param User body models.AdminLoginReq true "Login"
// @Success 201 {object} models.AdminLoginResp
// @Failure 400 string error models.Error
// @Failure 400 string error models.Error
// @Router /v1/auth/login [post]
func (h *handlerV1) LoginAdmin(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
		body        models.AdminLoginReq
	)
	jspbMarshal.UseProtoNames = true

	err := c.BindJSON(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", logger.Error(err))
		return
	}

	role, password, status, err := h.postgres.Get(body.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot get admin", logger.Error(err))
		return
	}

	if !status {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "this admin does not exist",
		})
		h.log.Error("admin does not exist", logger.Error(err))
		return
	}

	if !etc.CompareHashPassword(password, body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "incorrect password",
		})
		h.log.Error("incorrect password", logger.Error(err))
		return
	}

	if role == "admin" {
		h.jwtHandler = tokens.JWTHandler{
			Sub:       body.Username,
			Role:      "admin",
			SignInKey: h.cfg.SignInKey,
			Log:       h.log,
			Timeout:   h.cfg.AccessTokenTimeout,
		}
	} else if role == "superadmin" {
		fmt.Println(role)
		h.jwtHandler = tokens.JWTHandler{
			Sub:       body.Username,
			Role:      "superadmin",
			SignInKey: h.cfg.SignInKey,
			Log:       h.log,
			Timeout:   h.cfg.AccessTokenTimeout,
		}
	}

	access, _, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot create access token", logger.Error(err))
		return
	}

	response := models.AdminLoginResp{
		AccessToken: access,
	}

	c.JSON(http.StatusOK, response)
}
