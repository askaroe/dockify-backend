package location

import (
	"net/http"

	"github.com/askaroe/dockify-backend/internal/entity"
	"github.com/askaroe/dockify-backend/internal/services"
	"github.com/askaroe/dockify-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Location interface {
	GetNearestUsers(c *gin.Context)
}

type location struct {
	s      *services.Service
	logger *utils.Logger
}

func NewLocationHandler(s *services.Service, logger *utils.Logger) Location {
	return &location{s: s, logger: logger}
}

// GetNearestUsers godoc
// @Summary Get nearest users
// @Description Retrieves users nearest to given coordinates within a radius.
// @Tags location
// @Accept json
// @Produce json
// @Param request body entity.NearestUsersRequest true "Nearest users request"
// @Success 200 {array} entity.NearestUsersResponse
// @Success 204 {object} nil "no content"
// @Failure 400 {object} entity.ErrorMessage
// @Failure 500 {object} entity.ErrorMessage
// @Router /api/v1/location/nearest [post]
func (l *location) GetNearestUsers(c *gin.Context) {
	ctx := c.Request.Context()

	var request entity.NearestUsersRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "invalid request"})
		return
	}

	users, err := l.s.Location.GetNearestUsers(ctx, request)
	if err != nil {
		l.logger.Errorf("GetNearestUsers error: %v", err)
		c.JSON(http.StatusInternalServerError, entity.ErrorMessage{Message: "failed to get nearest users"})
		return
	}

	if len(users) == 0 || users == nil {
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, users)
}
