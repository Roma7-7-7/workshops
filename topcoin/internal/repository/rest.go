package rest

import (
	"encoding/json"
	"fmt"
	"github.com/Roma7-7-7/workshops/topcoin/internal/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Repository struct {
	client          *http.Client
	coinMarketToken string
}

func (r *Repository) GetCoinMarket(limit int) (*models.CoinMarketCap, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	req, err := http.NewRequest(http.MethodGet, "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?limit="+strconv.Itoa(limit), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-CMC_PRO_API_KEY", r.coinMarketToken)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			log.Printf("error while getting response from coinmarketcap: code=%v, body=%s", resp.StatusCode, string(body))
		} else {
			log.Printf("error while getting response from coinmarketcap: code=%v", resp.StatusCode)
		}
		return nil, fmt.Errorf("non 200 response code=%v", resp.StatusCode)
	}

	var result models.CoinMarketCap
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %v", err)
	}
	return &result, nil
}

func (r *Repository) GetCryptoCurrency(symbols ...string) (*models.CryptoCurrency, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("symbols is empty")
	}

	req, err := http.NewRequest(http.MethodGet, "https://min-api.cryptocompare.com/data/pricemulti?tsyms=USD&fsyms="+strings.Join(symbols, ","), nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			log.Printf("error while getting response from cryptocompare: code=%v, body=%s", resp.StatusCode, string(body))
		} else {
			log.Printf("error while getting response from cryptocompare: code=%v", resp.StatusCode)
		}
		return nil, fmt.Errorf("non 200 response code=%v", resp.StatusCode)
	}

	var result models.CryptoCurrency
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %v", err)
	}
	return &result, nil
}

func NewRepository(coinMarketToken string) Repository {
	return Repository{
		client:          &http.Client{},
		coinMarketToken: coinMarketToken,
	}
}
