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
// @Success 201 {object} entity.CreatedUserResponse "created user id"
// @Failure 400 {object} entity.ErrorMessage "invalid request"
// @Failure 500 {object} entity.ErrorMessage "failed to register user"
// @Router /api/v1/register [post]
func (u *user) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.logger.Errorf("Bind json error: %v", err)
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "invalid request"})
	}

	userId, err := u.s.User.Register(ctx, req)
	if err != nil {
		u.logger.Errorf("Register error: %v", err)
		c.JSON(http.StatusInternalServerError, entity.ErrorMessage{Message: "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, entity.CreatedUserResponse{UserID: userId})
}

// Login
// @Summary User login
// @Description Authenticate user and return user information
// @Tags User
// @Accept json
// @Produce json
// @Param request body entity.UserLoginRequest true "User login payload"
// @Success 200 {object} map[string]interface{} "authenticated user"
// @Failure 400 {object} entity.ErrorMessage "invalid request"
// @Failure 401 {object} entity.ErrorMessage "invalid email or password"
// @Failure 500 {object} entity.ErrorMessage "internal server error"
// @Router /api/v1/login [post]
func (u *user) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "invalid request"})
		return
	}

	userResponse, err := u.s.User.Login(ctx, req)
	if err != nil {
		u.logger.Errorf("Login error: %v", err)
		c.JSON(http.StatusUnauthorized, entity.ErrorMessage{Message: "failed to login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userResponse})
}
