package api

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserPassword struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Wallet struct {
	ID      string `json:"id"`
	UserId  string `json:"user_id"`
	Balance int    `json:"balance"`
}

type Transaction struct {
	ID             string `json:"id"`
	CreditWalletId string `json:"credit_wallet_id"`
	DebitWalletId  string `json:"debit_wallet_id"`
	Amount         int    `json:"amount"`
	Type           uint8  `json:"type"`
	FeeWalletId    string `json:"fee_wallet_id"`
	FeeAmount      int    `json:"fee_amount"`
}

type TransactionU struct {
	Transaction
	CreditUserID string `json:"credit_user_id"`
	DebitUserId  string `json:"debit_user_id"`
}
