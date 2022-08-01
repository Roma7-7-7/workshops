package calendar

import (
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
)

type Repository interface {
	GetEvents(title, dateFrom, timeFrom, dateTo, timeTo string) ([]models.Event, error)
}

// Service holds calendar business logic and works with repository
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetEvents(title, dateFrom, timeFrom, dateTo, timeTo string) ([]models.Event, error) {
	return s.repo.GetEvents(title, dateFrom, timeFrom, dateTo, timeTo)
}
