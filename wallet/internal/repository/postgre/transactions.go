package postgre

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Roma7-7-7/workshops/wallet/internal/models"
)

const FeeWalletId = "85aa7525-4fdb-4436-a600-66ffc55e0f65"

const defaultType = 1

func (r *Repository) GetTransactionsByUserID(userID string) ([]*models.UserTransaction, error) {
	query, args, err := psql.Select("t.id", "t.credit_wallet_id", "t.debit_wallet_id", "t.amount", "t.type", "t.fee_wallet_id", "cu.id", "du.id").
		From("transactions as t").
		LeftJoin("wallets cw on t.credit_wallet_id = cw.id").
		LeftJoin("wallets dw on t.debit_wallet_id = dw.id").
		LeftJoin("users cu on cu.id = cw.user_id").
		LeftJoin("users du on du.id = dw.user_id").
		Where(sq.Or{
			sq.Eq{"cu.id": userID},
			sq.Eq{"du.id": userID},
		}).
		OrderBy("t.id").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build transactions query: %v", err)
	}

	var rows *sql.Rows
	rows, err = r.db.Query(query, args...)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return nil, fmt.Errorf(`querying with sql="%s": %v`, query, err)
	}

	var transactions []*models.UserTransaction
	for rows.Next() {
		var t models.UserTransaction
		var amount uint64

		if err = rows.Scan(&t.ID, &t.CreditWalletID, &t.DebitWalletID, &amount, &t.Type, &t.FeeWalletID, &t.CreditUserID, &t.DebitUserID); err != nil {
			return nil, fmt.Errorf("scan wallet: %v", err)
		}
		t.Amount = models.AmountFromDBU(amount)
		transactions = append(transactions, &t)
	}

	return transactions, nil
}

func (r *Repository) TransferFunds(creditWalletID string, debitWalletId string, amount models.Amount, fee models.Amount) (*models.Transaction, error) {
	tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %v", err)
	}
	defer tx.Rollback()

	var result models.Transaction
	var dba uint64
	var dbFee uint64
	err = psql.Insert("transactions").
		Columns("credit_wallet_id", "debit_wallet_id", "amount", "type", "fee_amount", "fee_wallet_id").
		Values(creditWalletID, debitWalletId, amount.ToDB(), defaultType, fee.ToDB(), FeeWalletId).
		Suffix("RETURNING id, credit_wallet_id, debit_wallet_id, amount, type, fee_amount, fee_wallet_id").
		RunWith(tx).
		QueryRow().
		Scan(&result.ID, &result.CreditWalletID, &result.DebitWalletID, &dba, &result.Type, &dbFee, &result.FeeWalletID)
	if err != nil {
		return nil, fmt.Errorf("insert transaction: %v", err)
	}
	result.Amount = models.AmountFromDBU(dba)
	result.FeeAmount = models.AmountFromDBU(dbFee)

	if _, err = psql.Update("wallets").Set("balance", sq.Expr("balance + ?", amount.ToDB())).Where(sq.Eq{"id": debitWalletId}).RunWith(tx).Exec(); err != nil {
		return nil, fmt.Errorf("update debit wallet: %v", err)
	}
	if _, err = psql.Update("wallets").Set("balance", sq.Expr("balance + ?", fee.ToDB())).Where(sq.Eq{"id": FeeWalletId}).RunWith(tx).Exec(); err != nil {
		return nil, fmt.Errorf("update debit wallet: %v", err)
	}

	var creditB int64
	err = psql.Update("wallets").
		Set("balance", sq.Expr("balance - ? - ?", amount.ToDB(), fee.ToDB())).
		Where(sq.Eq{"id": creditWalletID}).
		Suffix("RETURNING balance").
		RunWith(tx).
		Scan(&creditB)
	if err != nil {
		return nil, fmt.Errorf("update credit wallet: %v", err)
	}
	if creditB < 0 {
		return nil, fmt.Errorf("insufficient funds")
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %v", err)
	}
	return &result, nil
}
