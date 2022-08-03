package calendar

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"time"
)

const DateTimeLayout = "2006-01-02 15:04"
const DateLayout = "2006-01-02"
const TimeLayout = "15:04"

type Repository interface {
	GetEvents(title, dateFrom, timeFrom, dateTo, timeTo string) ([]*models.Event, error)
	GetEvent(id string) (*models.Event, error)
	CreateEvent(title, description string, from time.Time, to time.Time, notes []string) (*models.Event, error)
	UpdateEvent(id, title, description string, from time.Time, to time.Time, notes []string) (*models.Event, error)
}

// Service holds calendar business logic and works with repository
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetEvents(title, dateFrom, timeFrom, dateTo, timeTo, timezone string) ([]*models.Event, error) {
	if timezone != "" {
		if dateFrom != "" {
			convertedDate, convertedTime, err := normalizeDateTime(dateFrom, timeFrom, timezone)
			if err != nil {
				return nil, fmt.Errorf("convert date=\"%s\" time=\"%s\" to timezone=\"%s\": %v", dateFrom, timeFrom, timezone, err)
			}
			dateFrom = convertedDate
			timeFrom = convertedTime
		}
		if dateTo != "" {
			convertedDate, convertedTime, err := normalizeDateTime(dateTo, timeTo, timezone)
			if err != nil {
				return nil, fmt.Errorf("convert date=\"%s\" time=\"%s\" to timezone=\"%s\": %v", dateTo, timeTo, timezone, err)
			}
			dateTo = convertedDate
			timeTo = convertedTime
		}
	}
	return s.repo.GetEvents(title, dateFrom, timeFrom, dateTo, timeTo)
}

func (s *Service) GetEvent(id string) (*models.Event, error) {
	return s.repo.GetEvent(id)
}

func (s *Service) CreateEvent(title, description, timeVal, timezone string, duration time.Duration, notes []string) (*models.Event, error) {
	timeFrom, timeTo, err := timeFromTo(timeVal, timezone, duration)
	if err != nil {
		return nil, err
	}
	return s.repo.CreateEvent(title, description, *timeFrom, *timeTo, notes)
}

func (s *Service) UpdateEvent(id, title, description, timeVal, timezone string, duration time.Duration, notes []string) (*models.Event, error) {
	timeFrom, timeTo, err := timeFromTo(timeVal, timezone, duration)
	if err != nil {
		return nil, err
	}
	return s.repo.UpdateEvent(id, title, description, *timeFrom, *timeTo, notes)
}

func timeFromTo(timeVal, timezone string, duration time.Duration) (*time.Time, *time.Time, error) {
	l, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid location %s: %v", timezone, err)
	}
	timeFrom, err := time.ParseInLocation(DateTimeLayout, timeVal, l)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid time %s: %v", timeVal, err)
	}
	timeTo := timeFrom.Add(duration)
	return &timeFrom, &timeTo, nil
}

func normalizeDateTime(date string, timev string, timezone string) (string, string, error) {
	if date == "" && timev == "" {
		return "", "", nil
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", "", fmt.Errorf("invalid timzone: %v", err)
	}

	if loc == nil {
		return date, timev, nil
	}

	dateTime := fmt.Sprintf("%s %s", date, timev)
	zoned, err := time.ParseInLocation(DateTimeLayout, dateTime, loc)
	if err != nil {
		return "", "", fmt.Errorf("convert date time %s: %v", dateTime, err)
	}
	converted := zoned.In(time.Now().Location())

	return converted.Format(DateLayout), converted.Format(TimeLayout), nil
}
