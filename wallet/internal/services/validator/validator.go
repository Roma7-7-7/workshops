package validator

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/wallet/internal/models"
	"go.uber.org/zap"
	"regexp"
	"strconv"
	"strings"
)

// Service validates structures
type Service struct {
}

type GetPageable struct {
	Limit  string
	Offset string
}

type CreateUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type CreateWallet struct {
	Balance string `json:"balance"`
}

type CreateTransaction struct {
	CreditWalletID string `json:"credit_wallet_id"`
	DebitWalletID  string `json:"debit_wallet_id"`
	Amount         string `json:"amount"`
}

func (s *Service) Validate(v interface{}) error {
	errors := make([]string, 0, 10)
	switch t := v.(type) {
	case *GetPageable:
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
	case *CreateWallet:
		if !IsValidAmount(t.Balance) {
			errors = append(errors, "balance must be a valid amount")
		}
	case *CreateTransaction:
		if strings.TrimSpace(t.CreditWalletID) == "" {
			errors = append(errors, "credit_wallet_id must not be empty")
		}
		if strings.TrimSpace(t.DebitWalletID) == "" {
			errors = append(errors, "debit_wallet_id must not be empty")
		}
		if t.CreditWalletID == t.DebitWalletID {
			errors = append(errors, "credit_wallet_id and debit_wallet_id must not be equal")
		}
		if !IsValidAmount(t.Amount) {
			errors = append(errors, "amount must be a valid amount")
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

func IsValidAmount(amount string) bool {
	if strings.TrimSpace(amount) == "" {
		return false
	}
	matches, err := regexp.MatchString(`^\d+(\.\d{1,2})?$`, amount)
	return matches && err == nil
}

func (u GetPageable) LimitN() uint64 {
	if strings.TrimSpace(u.Limit) == "" {
		return 20
	}
	// Ignore error as it should be handled by validation below
	res, _ := strconv.ParseUint(u.Limit, 10, 64)
	return res
}

func (u GetPageable) OffsetN() uint64 {
	if strings.TrimSpace(u.Offset) == "" {
		return 0
	}
	// Ignore error as it should be handled by validation below
	res, _ := strconv.ParseUint(u.Offset, 10, 64)
	return res
}

func (w CreateWallet) ToAmount() models.Amount {
	// Ignore error as it should be handled by validation below
	res, _ := models.ToAmount(w.Balance)
	return res
}
