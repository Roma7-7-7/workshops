package http

import "github.com/gin-gonic/gin"

// will hold http routes and will registrate them

func (s *Server) Register(e *gin.Engine) {
	e.POST("/login", s.auth.Login)
	e.GET("/logout", s.auth.Logout)

	api := e.Group("/api")
	api.Use(s.auth.Validate)
	s.registerEvents(api.Group("/events"))
}
