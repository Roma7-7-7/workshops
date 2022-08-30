package models

type User struct {
	ID       string
	Name     string
	Password string
}

type Wallet struct {
	ID      string
	Balance int64
	UserId  string
}

type Transaction struct {
	ID             string
	CreditWalletId string
	DebitWalletId  string
	Amount         int64
	Type           uint8
	FeeAmount      int64
	FeeWalletId    string
	CreditUserId   string
	DebitUserId    string
}
