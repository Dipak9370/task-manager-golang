package handler

import (
	"net/http"
	"task-manager/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service *service.AuthService
}

// Register godoc
// @Summary Register new user
// @Description Register user with email, password and role
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "User Register"
// @Success 200 {object} map[string]string
// @Router /register [post]

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.Register(c, req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registered"})
}

// Login godoc
// @Summary Login user
// @Description Login and get JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "User Login"
// @Success 200 {object} map[string]string
// @Router /login [post]

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	c.ShouldBindJSON(&req)

	token, err := h.Service.Login(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
