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
	Address string `json:"address"`
	UserId  string `json:"user_id"`
	Balance int    `json:"balance"`
}

type Transaction struct {
	ID            string `json:"id"`
	Amount        int    `json:"amount"`
	FeeAmount     int    `json:"fee_amount"`
	CreditUserID  string `json:"credit_user_id"`
	CreditAddress string `json:"credit_address"`
	DebitUserId   string `json:"debit_user_id"`
	DebitAddress  string `json:"debit_address"`
	FeeAddress    string `json:"fee_address"`
}
