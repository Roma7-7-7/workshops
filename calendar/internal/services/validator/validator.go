package validator

import "log"

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

func (s *Service) Validate(v interface{}) error {
	switch t := v.(type) {
	case GetEvents:
		return nil
	default:
		log.Printf("Validation of type %T is not supported\n", t)
		return nil
	}
}
