package models

type CoinMarketCap struct {
	Data []struct {
		Symbol string `json:"symbol"`
	} `json:"data"`
}

type CryptoCurrencyItem struct {
	USD float64 `json:"USD"`
}

type CryptoCurrency map[string]CryptoCurrencyItem
