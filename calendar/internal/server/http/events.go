package http

import (
	"github.com/Roma7-7-7/workshops/calendar/api"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/validator"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func (s *Server) GetEvents(c *gin.Context) {
	req := validator.GetEvents{
		Title:    c.Query("title"),
		DateFrom: c.Query("date_from"),
		TimeFrom: c.Query("time_from"),
		DateTo:   c.Query("date_to"),
		TimeTo:   c.Query("time_to"),
	}
	if err := s.valid.Validate(req); err != nil {
		badRequest(c, err)
		return
	}

	events, err := s.service.GetEvents(req.Title, req.DateFrom, req.TimeFrom, req.DateTo, req.TimeTo)
	if err != nil {
		log.Printf("get events: %v\n", err)
		serverError(c, err)
		return
	}

	result := make([]*api.Event, len(events))
	for i, e := range events {
		result[i] = toApi(e)
	}
	c.JSON(http.StatusOK, result)
}

func (s *Server) PostEvent(c *gin.Context) {
	var req validator.CreateEvent
	c.BindJSON(&req)
	if err := s.valid.Validate(req); err != nil {
		badRequest(c, err)
		return
	}

	e, err := s.service.CreateEvent(req.Title, req.Description, req.Time, req.Timezone, time.Duration(req.Duration)*time.Minute, req.Notes)
	if err != nil {
		log.Printf("create event: %v\n", err)
		serverError(c, err)
		return
	}

	c.JSON(http.StatusOK, toApi(e))
}

func toApi(e *models.Event) *api.Event {
	var tz string
	if l := e.TimeFrom.Location(); l == nil {
		tz = "UTC"
	} else {
		tz = l.String()
	}
	return &api.Event{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Time:        e.TimeFrom.Format("2006-01-02 15:04"),
		TimeZone:    tz,
		Duration:    int(e.TimeTo.Sub(e.TimeFrom).Minutes()),
		Notes:       e.Notes,
	}
}
