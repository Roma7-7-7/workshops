package wallet

import (
	"errors"
	"fmt"
	"github.com/Roma7-7-7/workshops/wallet/internal/models"
	"github.com/Roma7-7-7/workshops/wallet/internal/repository/postgre"
)

type Repository interface {
	GetUserByName(name string) (*models.User, error)
	GetUsers(limit uint64, offset uint64) ([]*models.User, error)
	CreateUser(name string, password string) (*models.User, error)

	CreateWallet(userID string, balance models.Amount) (*models.Wallet, error)
	GetWalletOwner(id string) (string, error)
	GetWalletByID(id string) (*models.Wallet, error)
	GetWalletTransactionsU(id string) (*models.Wallet, []*models.UserTransaction, error)

	GetTransactionsByUserID(userID string) ([]*models.UserTransaction, error)
	TransferFunds(creditWalletID string, debitWalletId string, amount models.Amount, fee models.Amount) (*models.Transaction, error)
}

const fee = 1.5

// Service holds wallet business logic and works with repository
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUsers(limit uint64, offset uint64) ([]*models.User, error) {
	return s.repo.GetUsers(limit, offset)
}

var ErrUserExists = errors.New("user already exists")

func (s *Service) CreateUser(name string, password string) (*models.User, error) {
	if u, err := s.repo.GetUserByName(name); err != nil {
		return nil, fmt.Errorf("get user by name: %v", err)
	} else if u != nil {
		return nil, ErrUserExists
	}

	return s.repo.CreateUser(name, password)
}

func (s *Service) CreateWallet(userID string, balance models.Amount) (*models.Wallet, error) {
	return s.repo.CreateWallet(userID, balance)
}

func (s *Service) GetWalletByID(id string) (*models.Wallet, error) {
	return s.repo.GetWalletByID(id)
}

func (s *Service) GetWalletTransactionsU(id string) (*models.Wallet, []*models.UserTransaction, error) {
	return s.repo.GetWalletTransactionsU(id)
}

func (s *Service) GetWalletOwner(id string) (string, error) {
	return s.repo.GetWalletOwner(id)
}

func (s *Service) GetTransactionsByUserID(userID string) ([]*models.UserTransaction, error) {
	return s.repo.GetTransactionsByUserID(userID)
}

var ErrFeeWalletTransfer = errors.New("transfer directly from/to fee wallet is not allowed")

func (s *Service) TransferFunds(creditWalletID string, debitWalletID string, amount models.Amount) (*models.Transaction, error) {
	if creditWalletID == debitWalletID {
		return nil, errors.New("wallets are equal")
	}
	if creditWalletID == postgre.FeeWalletId || debitWalletID == postgre.FeeWalletId {
		return nil, ErrFeeWalletTransfer
	}

	return s.repo.TransferFunds(creditWalletID, debitWalletID, amount, models.Amount(uint64(float64(amount)/100*fee)))
}
