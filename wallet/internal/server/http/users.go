package http

import (
	"github.com/Roma7-7-7/workshops/wallet/api"
	"github.com/Roma7-7-7/workshops/wallet/internal/models"
	"github.com/Roma7-7-7/workshops/wallet/internal/services/validator"
	"github.com/Roma7-7-7/workshops/wallet/internal/services/wallet"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (s *Server) GetUsers(c *gin.Context) {
	req := validator.GetPageable{
		Limit:  c.Query("limit"),
		Offset: c.Query("offset"),
	}
	zap.L().Debug("get users", zap.Any("payload", req))
	if err := s.valid.Validate(&req); err != nil {
		api.BadRequestA(c, err)
		return
	}

	users, err := s.service.GetUsers(req.LimitN(), req.OffsetN())
	if err != nil {
		zap.L().Error("get users", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	}

	result := make([]*api.User, len(users))
	for i, e := range users {
		result[i] = userToApi(e)
	}
	c.JSON(http.StatusOK, result)
}

func (s *Server) PostUser(c *gin.Context) {
	req := &validator.CreateUser{}
	if err := c.ShouldBindJSON(req); err != nil {
		api.BadJSONA(c)
		return
	}

	zap.L().Debug("put user", zap.Any("payload", req))
	if err := s.valid.Validate(req); err != nil {
		api.BadRequestA(c, err)
		return
	}

	user, err := s.service.CreateUser(req.Name, req.Password)
	if err == wallet.ErrUserExists {
		zap.L().Debug("user already exists", zap.Error(err))
		api.ConflictA(c, "user already exists")
		return
	} else if err != nil {
		zap.L().Error("put user", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	}

	c.JSON(http.StatusOK, userToApi(user))
}

func userToApi(u *models.User) *api.User {
	return &api.User{
		ID:   u.ID,
		Name: u.Name,
	}
}
