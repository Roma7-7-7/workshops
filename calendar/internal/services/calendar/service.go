package calendar

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"time"
)

const TimeLayout = "2006-01-02 03:04"

type Repository interface {
	GetEvents(title, dateFrom, timeFrom, dateTo, timeTo string) ([]*models.Event, error)
	CreateEvent(title string, description string, from time.Time, to time.Time, notes []string) (*models.Event, error)
}

// Service holds calendar business logic and works with repository
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetEvents(title, dateFrom, timeFrom, dateTo, timeTo string) ([]*models.Event, error) {
	return s.repo.GetEvents(title, dateFrom, timeFrom, dateTo, timeTo)
}

func (s *Service) CreateEvent(title, description, timeVal, timezone string, duration time.Duration, notes []string) (*models.Event, error) {
	l, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid location %s: %v", timezone, err)
	}
	timeFrom, err := time.ParseInLocation(TimeLayout, timeVal, l)
	if err != nil {
		return nil, fmt.Errorf("invalid time %s: %v", timeVal, err)
	}

	return s.repo.CreateEvent(title, description, timeFrom, timeFrom.Add(duration), notes)
}
