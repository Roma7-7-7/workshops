package http

import (
	"github.com/Roma7-7-7/workshops/calendar/internal/middleware/auth"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/calendar"
)

type Validator interface {
	Validate(interface{}) error
}

type Server struct {
	service *calendar.Service
	auth    *auth.Middleware
	valid   Validator
	secret  string
}

func NewServer(service *calendar.Service, auth *auth.Middleware, valid Validator, secret string) *Server {
	return &Server{
		service: service,
		auth:    auth,
		valid:   valid,
		secret:  secret,
	}
}
