package postgre

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Roma7-7-7/workshops/wallet/internal/models"
	_ "github.com/lib/pq"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	db *sql.DB
}

func NewRepository(dsn string) *Repository {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return &Repository{
		db: db,
	}
}

func (r *Repository) GetUserByName(name string) (*models.User, error) {
	var user models.User
	err := psql.Select("id", "name", "password").
		From("users").
		Where(sq.Eq{"name": name}).
		RunWith(r.db).
		QueryRow().
		Scan(&user.ID, &user.Name, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("get user: %v", err)
	}
	return &user, nil
}
