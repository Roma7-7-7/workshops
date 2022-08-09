package http

import (
	"github.com/Roma7-7-7/workshops/calendar/api"
	"github.com/Roma7-7-7/workshops/calendar/internal/middleware/auth"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) UpdateUserTimezone(c *gin.Context) {
	var req api.UserTimezone
	if err := c.BindJSON(&req); err != nil {
		api.BadRequestA(c, err)
		return
	}

	if err := s.valid.Validate(&req); err != nil {
		api.BadRequestA(c, err)
		return
	}

	user := auth.GetContext(c).Username()
	if user != req.Username {
		api.ForbiddenA(c, "You are not allowed to update this user")
		return
	}

	u, err := s.service.UpdateUserTimezone(req.Username, req.Timezone)
	if err != nil {
		api.ServerErrorA(c, err)
		return
	}

	c.JSON(http.StatusOK, userToApi(u))
}

func userToApi(u *models.User) *api.UserTimezone {
	return &api.UserTimezone{
		Username: u.Name,
		Timezone: u.Timezone,
	}
}