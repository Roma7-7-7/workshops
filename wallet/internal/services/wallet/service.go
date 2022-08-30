package wallet

import (
	"github.com/Roma7-7-7/workshops/wallet/internal/models"
)

type Repository interface {
	GetUserByName(name string) (*models.User, error)
}

// Service holds wallet business logic and works with repository
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
