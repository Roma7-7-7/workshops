package http

import (
	"errors"
	"github.com/Roma7-7-7/workshops/wallet/api"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Validator interface {
	Validate(interface{}) error
}

type Service interface {
}

type Auth interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	ValidateGin(c *gin.Context)
}

type Server struct {
	service Service
	valid   Validator
	auth    Auth
}

func NewServer(service Service, valid Validator, auth Auth) *Server {
	return &Server{
		service: service,
		valid:   valid,
		auth:    auth,
	}
}

func (s *Server) Register(app *gin.Engine) {
	app.Use(gin.CustomRecovery(recoveryHandler))

	app.POST("/login", s.auth.Login)
	app.GET("/logout", s.auth.Logout)

}

func recoveryHandler(c *gin.Context, err interface{}) {
	zap.L().Error("unexpected panic", zap.Any("panic", err))
	api.ServerErrorA(c, errors.New("unexpected error occurred"))
}
