package api

type Coin struct {
	Rank     int      `json:"rank"`
	Symbol   string   `json:"symbol"`
	PriceUSD *float64 `json:"price_usd"`
}
