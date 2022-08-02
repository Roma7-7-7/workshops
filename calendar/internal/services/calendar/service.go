package calendar

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"time"
)

const DateTimeLayout = "2006-01-02 03:04"
const DateLayout = "2006-01-02"
const TimeLayout = "15:04"

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

func (s *Service) CreateEvent(title, description, timeVal, timezone string, duration time.Duration, notes []string) (*models.Event, error) {
	l, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid location %s: %v", timezone, err)
	}
	timeFrom, err := time.ParseInLocation(DateTimeLayout, timeVal, l)
	if err != nil {
		return nil, fmt.Errorf("invalid time %s: %v", timeVal, err)
	}

	return s.repo.CreateEvent(title, description, timeFrom, timeFrom.Add(duration), notes)
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
