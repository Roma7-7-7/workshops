package http

import "github.com/gin-gonic/gin"

// will hold http routes and will registrate them

func (s *Server) Register(e *gin.Engine) {
	api := e.Group("/api")
	api.GET("/events", s.GetEvents)
	api.POST("/events", s.PostEvent)
}
