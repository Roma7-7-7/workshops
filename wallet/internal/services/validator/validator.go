package validator

import (
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

// Service validates structures
type Service struct {
}

type GetUsers struct {
	Limit  string
	Offset string
}

type CreateUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u GetUsers) LimitN() uint64 {
	if strings.TrimSpace(u.Limit) == "" {
		return 20
	}
	// Ignore error as it should be handled by validation below
	res, _ := strconv.ParseUint(u.Limit, 10, 64)
	return res
}

func (u GetUsers) OffsetN() uint64 {
	if strings.TrimSpace(u.Offset) == "" {
		return 0
	}
	// Ignore error as it should be handled by validation below
	res, _ := strconv.ParseUint(u.Offset, 10, 64)
	return res
}

func (s *Service) Validate(v interface{}) error {
	errors := make([]string, 0, 10)
	switch t := v.(type) {
	case *GetUsers:
		if val, err := strconv.Atoi(t.Limit); strings.TrimSpace(t.Limit) != "" && (err != nil || val < 0) {
			errors = append(errors, "limit must be a positive integer")
		}
		if val, err := strconv.Atoi(t.Offset); strings.TrimSpace(t.Offset) != "" && (err != nil || val < 0) {
			errors = append(errors, "offset must be a positive integer")
		}
	case *CreateUser:
		if strings.TrimSpace(t.Name) == "" {
			errors = append(errors, "username must not be empty")
		}
		if strings.TrimSpace(t.Password) == "" {
			errors = append(errors, "password must not be empty")
		}
	default:
		zap.L().Warn("validation is not supported", zap.Any("type", t))
		return nil
	}
	if len(errors) > 0 {
		return fmt.Errorf("validation errors: [%s]", strings.Join(errors, "; "))
	}
	return nil
}
