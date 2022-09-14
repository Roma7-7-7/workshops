package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type User struct {
	ID       string
	Name     string
	Password string
}

type Amount uint64

type Wallet struct {
	ID      string
	Balance Amount
	UserID  string
}

type Transaction struct {
	ID             string
	CreditWalletID string
	DebitWalletID  string
	Amount         Amount
	Type           uint8
	FeeAmount      Amount
	FeeWalletID    string
}

type UserTransaction struct {
	Transaction
	CreditUserID string
	DebitUserID  string
}

func AmountFromDB(value int64) Amount {
	return Amount(value)
}

func AmountFromDBU(value uint64) Amount {
	return Amount(value)
}

func ToAmount(value string) (Amount, error) {
	if strings.TrimSpace(value) == "" {
		return Amount(0), nil
	}
	if matches, err := regexp.MatchString(`^\d+(\.\d{1,2})?$`, value); !matches || err != nil {
		return Amount(0), fmt.Errorf("invalid amount %s", value)
	}
	res, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return Amount(0), fmt.Errorf("invalid amount %s: %v", value, err)
	}

	return Amount(res * 100), nil
}

func (a Amount) RoundWholePart() int {
	return int(a / 100)
}

func (a Amount) ToDB() uint64 {
	return uint64(a)
}
