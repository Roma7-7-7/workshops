package http

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/wallet/api"
	"github.com/Roma7-7-7/workshops/wallet/internal/middleware/auth"
	"github.com/Roma7-7-7/workshops/wallet/internal/models"
	"github.com/Roma7-7-7/workshops/wallet/internal/services/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (s *Server) PostWallet(c *gin.Context) {
	var req validator.CreateWallet
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	zap.L().Debug("post wallet", zap.Any("payload", req))
	if err := s.valid.Validate(&req); err != nil {
		api.BadRequestA(c, err)
		return
	}

	created, err := s.service.CreateWallet(auth.GetContext(c).UserID(), req.ToAmount())
	if err != nil {
		zap.L().Error("create wallet", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	}

	c.JSON(http.StatusOK, walletToApi(created))
}

func (s *Server) GetWallet(c *gin.Context) {
	id := c.Param("id")
	zap.L().Debug("get wallet", zap.String("id", id))
	if !s.validateOwner(c, id) {
		return
	}

	event, err := s.service.GetWalletByID(id)
	if err != nil {
		zap.L().Error("get wallet", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	}

	c.JSON(http.StatusOK, walletToApi(event))
}

func (s *Server) GetWalletWithUserTransactions(c *gin.Context) {
	id := c.Param("id")
	zap.L().Debug("get wallet with user transactions", zap.String("id", id))
	if !s.validateOwner(c, id) {
		return
	}

	_, transactions, err := s.service.GetWalletTransactionsU(id)
	if err != nil {
		zap.L().Error("get wallet with user transactions", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	}

	result := make([]*api.TransactionU, len(transactions))
	for i, t := range transactions {
		result[i] = transactionUToApi(t)
	}
	c.JSON(http.StatusOK, result)
}

func (s *Server) registerWallet(e *gin.Engine, middleware ...gin.HandlerFunc) {
	e.POST("/wallets", append(middleware, s.PostWallet)...)
	e.GET("/wallets/:id", append(middleware, s.GetWallet)...)
	e.GET("/wallets/:id/transactions", append(middleware, s.GetWalletWithUserTransactions)...)
}

func (s *Server) validateOwner(c *gin.Context, walletID string) bool {
	if owner, err := s.service.GetWalletOwner(walletID); err != nil {
		zap.L().Error("get event owner of wallet", zap.String("ID", walletID), zap.Error(err))
		api.ServerErrorA(c, err)
		return false
	} else if owner == "" {
		api.NotFoundA(c, fmt.Sprintf("wallet with ID=\"%s\"", walletID))
		return false
	} else if owner != auth.GetContext(c).UserID() {
		api.ForbiddenA(c, fmt.Sprintf("wallet with ID=\"%s\"", walletID))
		return false
	}
	return true
}

func walletToApi(w *models.Wallet) *api.Wallet {
	return &api.Wallet{
		ID:      w.ID,
		UserID:  w.UserID,
		Balance: w.Balance.RoundWholePart(),
	}
}
