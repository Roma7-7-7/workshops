package postgre

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Roma7-7-7/workshops/wallet/internal/models"
)

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

func (r *Repository) GetUsers(limit uint64, offset uint64) ([]*models.User, error) {
	q := psql.Select("id", "name", "password").
		From("users").
		Limit(limit).
		Offset(offset).
		OrderBy("name")

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build users query: %v", err)
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

	var result []*models.User
	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Password); err != nil {
			return nil, fmt.Errorf("scan user: %v", err)
		}
		result = append(result, &user)
	}

	return result, nil
}

func (r *Repository) CreateUser(name string, password string) (*models.User, error) {
	var user models.User
	err := psql.Insert("users").
		Columns("name", "password").
		Values(name, password).
		Suffix("RETURNING id, name, password").
		RunWith(r.db).
		QueryRow().
		Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("create user: %v", err)
	}

	return &user, nil
}
