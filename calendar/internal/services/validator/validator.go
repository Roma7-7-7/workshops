package validator

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/calendar"
	"log"
	"strings"
	"time"
)

// Service validates structures
type Service struct {
}

type GetEvents struct {
	Title    string
	DateFrom string
	TimeFrom string
	DateTo   string
	TimeTo   string
}

type CreateEvent struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Time        string   `json:"time"`
	Timezone    string   `json:"timezone"`
	Duration    int      `json:"duration"`
	Notes       []string `json:"notes"`
}

func (s *Service) Validate(v interface{}) error {
	errors := make([]string, 0, 10)
	switch t := v.(type) {
	case GetEvents:
		return nil
	case CreateEvent:
		if strings.TrimSpace(t.Title) == "" {
			errors = append(errors, "title must not be blank")
		}
		if _, err := time.Parse(calendar.TimeLayout, t.Time); err != nil {
			errors = append(errors, "invalid time value")
		}
		if _, err := time.LoadLocation(t.Timezone); err != nil || strings.TrimSpace(t.Timezone) == "" {
			errors = append(errors, "timezone must not be blank and be a valid timezone")
		}
		if t.Duration <= 0 {
			errors = append(errors, "duration must be greater than 0")
		}
		if len(errors) > 0 {
			return fmt.Errorf("validation errors: [%s]", strings.Join(errors, "; "))
		}
		return nil
	default:
		log.Printf("Validation of type %T is not supported\n", t)
		return nil
	}
}
