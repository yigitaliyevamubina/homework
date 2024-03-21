package v1

import (
	"exam/api-gateway/api/handlers/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// List all roles
// @Router /v1/rbac/roles [get]
// @Security BearerAuth
// @Summary get all roles
// @Tags Role-management
// @Description Get all roles
// @Accept json
// @Product json
// @Param super-username query string true "super-username"
// @Param super-password query string true "super-password"
// @Success 201 {object} models.RbacAllRolesResp
// @Failure 400 string error models.Error
// @Failure 400 string error models.Error
func (h *handlerV1) ListAllRoles(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	superAdminUsername := c.Query("super-username")
	superAdminPassword := c.Query("super-password")
	if superAdminPassword == "b" && superAdminUsername == "a" {
		roles := h.casbin.GetAllRoles()
		c.JSON(http.StatusOK, models.RbacAllRolesResp{
			Roles: roles,
		})
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "you cannot get all roles, provide correct username and password",
		})
	}
}

// List all policies of a role
// @Router /v1/rbac/policies/{role} [get]
// @Security BearerAuth
// @Summary get all policies of a role
// @Tags Role-management
// @Description Get all policies of a role
// @Accept json
// @Product json
// @Param super-username query string true "super-username"
// @Param super-password query string true "super-password"
// @Param role path string true "role"
// @Success 201 {object} models.ListRolePolicyResp
// @Failure 400 string error models.Error
// @Failure 400 string error models.Error
func (h *handlerV1) ListRolePolicies(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true
	superAdminUsername := c.Query("super-username")
	superAdminPassword := c.Query("super-password")
	if superAdminPassword == "b" && superAdminUsername == "a" {
		role := c.Param("role")
		var response models.ListRolePolicyResp
		for _, p := range h.casbin.GetFilteredPolicy(0, role) {
			response.Policies = append(response.Policies, &models.Policy{
				Role:     p[0],
				EndPoint: p[1],
				Method:   p[2],
			})
		}
		c.JSON(http.StatusOK, response)
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "you cannot get role policies, provide correct username and password",
		})
	}
}

// Add policy to a role
// @Router /v1/rbac/add/policy [post]
// @Security BearerAuth
// @Summary add policy to a role
// @Tags Role-management
// @Description Add policy to a role
// @Accept json
// @Product json
// @Param super-username query string true "super-username"
// @Param super-password query string true "super-password"
// @Param policy body models.AddPolicyRequest true "policy"
// @Success 201 {object} models.SuperAdminMessage
// @Failure 400 string error models.Error
// @Failure 400 string error models.Error
func (h *handlerV1) AddPolicyToRole(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
		body        models.AddPolicyRequest
	)
	jspbMarshal.UseProtoNames = true

	superAdminUsername := c.Query("super-username")
	superAdminPassword := c.Query("super-password")
	if superAdminPassword == "b" && superAdminUsername == "a" {
		err := c.BindJSON(&body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "failed, try again",
			})
			return
		}
		p := []string{body.Policy.Role, body.Policy.EndPoint, body.Policy.Method}
		if _, err := h.casbin.AddPolicy(p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "failed, try again",
			})
			return
		}
		h.casbin.SavePolicy()
		c.JSON(http.StatusOK, models.SuperAdminMessage{
			Message: "success",
		})
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "you cannot grand permission to the admin role, provide correct username and password",
		})
	}
}

// Delete policy
// @Router /v1/rbac/delete/policy [delete]
// @Security BearerAuth
// @Summary delete policy
// @Tags Role-management
// @Description Delete policy
// @Accept json
// @Product json
// @Param super-username query string true "super-username"
// @Param super-password query string true "super-password"
// @Param policy body models.AddPolicyRequest true "policy"
// @Success 201 {object} models.SuperAdminMessage
// @Failure 400 string error models.Error
// @Failure 400 string error models.Error
func (h *handlerV1) DeletePolicy(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
		body        models.AddPolicyRequest
	)
	jspbMarshal.UseProtoNames = true
	superAdminUsername := c.Query("super-username")
	superAdminPassword := c.Query("super-password")
	if superAdminPassword == "b" && superAdminUsername == "a" {
		err := c.BindJSON(&body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error,
				"message": "failed, try again",
			})
			return
		}
		p := []string{body.Policy.Role, body.Policy.EndPoint, body.Policy.Method}
		if _, err := h.casbin.RemovePolicy(p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error,
				"message": "failed, try again",
			})
			return
		}
		h.casbin.SavePolicy()
		c.JSON(http.StatusOK, models.SuperAdminMessage{
			Message: "success",
		})
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "you cannot grand permission to the admin role, provide correct username and password",
		})
	}
}
