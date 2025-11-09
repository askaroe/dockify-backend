package user

import (
	"net/http"

	"github.com/askaroe/dockify-backend/internal/entity"
	"github.com/askaroe/dockify-backend/internal/services"
	"github.com/askaroe/dockify-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type User interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type user struct {
	s      *services.Service
	logger *utils.Logger
}

func NewUserHandler(s *services.Service, logger *utils.Logger) User {
	return &user{
		s:      s,
		logger: logger}
}

// Register
// @Summary Register a new user
// @Description Create a new user account
// @Tags User
// @Accept json
// @Produce json
// @Param request body entity.UserRegisterRequest true "User registration payload"
// @Success 201 {object} map[string]interface{} "created user id"
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 500 {object} map[string]string "failed to register user"
// @Router /api/v1/register [post]
func (u *user) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.logger.Errorf("Bind json error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}

	userId, err := u.s.User.Register(ctx, req)
	if err != nil {
		u.logger.Errorf("Register error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": userId})
}

// Login
// @Summary User login
// @Description Authenticate user and return user information
// @Tags User
// @Accept json
// @Produce json
// @Param request body entity.UserLoginRequest true "User login payload"
// @Success 200 {object} map[string]interface{} "authenticated user"
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 401 {object} map[string]string "invalid email or password"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /api/v1/login [post]
func (u *user) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userResponse, err := u.s.User.Login(ctx, req)
	if err != nil {
		u.logger.Errorf("Login error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userResponse})
}
