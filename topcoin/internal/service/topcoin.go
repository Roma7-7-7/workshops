package service

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/topcoin/api"
	"github.com/Roma7-7-7/workshops/topcoin/internal/models"
)

type Repository interface {
	GetCoinMarket(limit int) (*models.CoinMarketCap, error)
	GetCryptoCurrency(symbols ...string) (*models.CryptoCurrency, error)
}

type Service struct {
	api Repository
}

func (s *Service) GetTopCoin(limit int) ([]*api.Coin, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	coinMarket, err := s.api.GetCoinMarket(limit)
	if err != nil {
		return nil, fmt.Errorf("getting coinmarket data: %v", err)
	}

	var symbols []string
	for _, coin := range coinMarket.Data {
		symbols = append(symbols, coin.Symbol)
	}

	cryptoCurrency, err := s.api.GetCryptoCurrency(symbols...)
	if err != nil {
		return nil, fmt.Errorf("getting crypto currency data: %v", err)
	}

	return merge(coinMarket, cryptoCurrency), nil
}
func NewService(api Repository) *Service {
	return &Service{api: api}
}

func merge(market *models.CoinMarketCap, currency *models.CryptoCurrency) []*api.Coin {
	var result []*api.Coin
	for i, coin := range market.Data {
		coin := &api.Coin{
			Rank:   i + 1,
			Symbol: coin.Symbol,
		}
		if price, ok := (*currency)[coin.Symbol]; ok {
			coin.PriceUSD = &price.USD
		}
		result = append(result, coin)
	}
	return result
}
