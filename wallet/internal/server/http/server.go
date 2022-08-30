package http

import (
	"github.com/Roma7-7-7/workshops/wallet/internal/models"
	"github.com/gin-gonic/gin"
)

type Validator interface {
	Validate(interface{}) error
}

type Service interface {
	GetUsers(limit uint64, offset uint64) ([]*models.User, error)
	CreateUser(name string, password string) (*models.User, error)
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
