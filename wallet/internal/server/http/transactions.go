package http

import (
	"github.com/Roma7-7-7/workshops/wallet/api"
	"github.com/Roma7-7-7/workshops/wallet/internal/middleware/auth"
	"github.com/Roma7-7-7/workshops/wallet/internal/models"
	"github.com/Roma7-7-7/workshops/wallet/internal/services/validator"
	"github.com/Roma7-7-7/workshops/wallet/internal/services/wallet"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (s *Server) GetTransactions(c *gin.Context) {
	req := validator.GetPageable{
		Limit:  c.Query("limit"),
		Offset: c.Query("offset"),
	}
	zap.L().Debug("get transactions", zap.Any("payload", req))
	if err := s.valid.Validate(&req); err != nil {
		api.BadRequestA(c, err)
		return
	}

	transactions, err := s.service.GetTransactionsByUserID(auth.GetContext(c).UserID())
	if err != nil {
		zap.L().Error("get user transactions", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	}

	result := make([]*api.TransactionU, len(transactions))
	for i, t := range transactions {
		result[i] = transactionUToApi(t)
	}
	c.JSON(http.StatusOK, result)
}

func (s *Server) PutTransaction(c *gin.Context) {
	req := &validator.CreateTransaction{}
	if err := c.ShouldBindJSON(req); err != nil {
		api.BadRequestA(c, err)
		return
	}

	zap.L().Debug("create transaction", zap.Any("payload", req))
	if err := s.valid.Validate(req); err != nil {
		api.BadRequestA(c, err)
		return
	}

	if !s.validateOwner(c, req.CreditWalletID) {
		return
	}

	// Ignore err because of above validation
	amount, _ := models.ToAmount(req.Amount)
	res, err := s.service.TransferFunds(req.CreditWalletID, req.DebitWalletID, amount)
	if err == wallet.ErrFeeWalletTransfer {
		api.BadRequestA(c, err)
		return
	} else if err != nil {
		zap.L().Error("create transaction", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	}

	c.JSON(http.StatusOK, transactionToApi(res))
}

func (s *Server) registerTransactions(e *gin.Engine, middleware ...gin.HandlerFunc) {
	e.GET("/transactions", append(middleware, s.GetTransactions)...)
	e.PUT("/transactions", append(middleware, s.PutTransaction)...)
}

func transactionToApi(t *models.Transaction) *api.Transaction {
	return &api.Transaction{
		ID:             t.ID,
		CreditWalletID: t.CreditWalletID,
		DebitWalletID:  t.DebitWalletID,
		Amount:         t.Amount.RoundWholePart(),
		Type:           t.Type,
		FeeWalletID:    t.FeeWalletID,
		FeeAmount:      t.FeeAmount.RoundWholePart(),
	}
}

func transactionUToApi(t *models.UserTransaction) *api.TransactionU {
	return &api.TransactionU{
		Transaction:  *transactionToApi(&t.Transaction),
		CreditUserID: t.CreditUserID,
		DebitUserID:  t.DebitUserID,
	}
}
