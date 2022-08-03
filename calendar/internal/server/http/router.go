package http

import "github.com/gin-gonic/gin"

// will hold http routes and will registrate them

func (s *Server) Register(e *gin.Engine) {
	api := e.Group("/api")
	events := api.Group("/events")
	events.GET("/", s.GetEvents)
	events.GET("/:id", s.GetEvent)
	events.POST("/", s.PostEvent)
	events.PUT("/:id", s.PutEvent)
}
