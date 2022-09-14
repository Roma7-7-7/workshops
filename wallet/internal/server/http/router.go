package http

import (
	"errors"
	"github.com/Roma7-7-7/workshops/wallet/api"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// will hold http routes and will registrate them

func (s *Server) Register(app *gin.Engine) {
	app.Use(gin.CustomRecovery(recoveryHandler))

	app.POST("/login", s.auth.Login)
	app.GET("/logout", s.auth.Logout)

	// No auth for users API
	app.GET("/users", s.GetUsers)
	app.POST("/users", s.PostUser)

	s.registerWallet(app, s.auth.ValidateGin)
	s.registerTransactions(app, s.auth.ValidateGin)
}

func recoveryHandler(c *gin.Context, err interface{}) {
	zap.L().Error("unexpected panic", zap.Any("panic", err))
	api.ServerErrorA(c, errors.New("unexpected error occurred"))
}
