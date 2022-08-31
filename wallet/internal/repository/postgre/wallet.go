package postgre

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Roma7-7-7/workshops/wallet/internal/models"
)

func (r *Repository) CreateWallet(userId string, balance models.Amount) (*models.Wallet, error) {
	var wallet models.Wallet
	err := psql.Insert("wallets").
		Columns("balance", "user_id").
		Values(balance, userId).
		Suffix("RETURNING id, balance, user_id").
		RunWith(r.db).
		QueryRow().
		Scan(&wallet.ID, &wallet.Balance, &wallet.UserId)

	if err != nil {
		return nil, fmt.Errorf("create wallet: %v", err)
	}

	return &wallet, nil
}

func (r *Repository) GetWalletOwner(id string) (string, error) {
	var userId string
	err := psql.Select("user_id").
		From("wallets").
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		QueryRow().
		Scan(&userId)

	if err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("get wallet owner: %v", err)
	}

	return userId, nil
}

func (r *Repository) GetWalletById(id string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := psql.Select("id", "balance", "user_id").
		From("wallets").
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		QueryRow().
		Scan(&wallet.ID, &wallet.Balance, &wallet.UserId)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("get wallet: %v", err)
	}

	return &wallet, nil
}

func (r *Repository) GetWalletTransactionsU(id string) (*models.Wallet, []*models.UserTransaction, error) {
	query, args, err := psql.Select("w.id", "w.balance", "w.user_id",
		"t.id", "t.credit_wallet_id", "t.debit_wallet_id", "t.amount", "t.type", "t.fee_wallet_id",
		"cu.id", "du.id").
		From("wallets as w").
		LeftJoin("transactions as t on w.id = t.credit_wallet_id or w.id = t.debit_wallet_id").
		LeftJoin("users cu on cu.id = t.credit_wallet_id").
		LeftJoin("users du on du.id = t.debit_wallet_id").
		Where(sq.Eq{"w.id": id}).
		RunWith(r.db).
		ToSql()

	if err != nil {
		return nil, nil, fmt.Errorf("build wallet query: %v", err)
	}

	var rows *sql.Rows
	rows, err = r.db.Query(query, args...)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return nil, nil, fmt.Errorf(`querying with sql="%s": %v`, query, err)
	}

	var wallet models.Wallet
	var transactions []*models.UserTransaction
	for rows.Next() {
		var tId, cwId, dwId, fwId, cuId, duId sql.NullString
		var tType, tAmount sql.NullInt64

		if err = rows.Scan(&wallet.ID, &wallet.Balance, &wallet.UserId,
			&tId, &cwId, &dwId, &tAmount, &tType, &fwId,
			&cuId, &duId); err != nil {
			return nil, nil, fmt.Errorf("scan wallet: %v", err)
		}
		if tId.Valid {
			transactions = append(transactions, &models.UserTransaction{
				Transaction: models.Transaction{
					ID:             tId.String,
					CreditWalletId: cwId.String,
					DebitWalletId:  dwId.String,
					Amount:         models.AmountFromDB(tAmount.Int64),
					Type:           uint8(tType.Int64),
					FeeWalletId:    fwId.String,
				},
				CreditUserID: cuId.String,
				DebitUserID:  duId.String,
			})
		}
	}

	return &wallet, transactions, nil
}
