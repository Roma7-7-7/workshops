package postgre

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
)

func (r *Repository) GetUser(name string) (*models.User, error) {
	var user models.User
	err := psql.Select("name", "password", "timezone").
		From("users").
		Where(sq.Eq{"name": name}).
		RunWith(r.db).
		QueryRow().
		Scan(&user.Name, &user.Password, &user.Timezone)
	if err != nil {
		return nil, fmt.Errorf("get user: %v", err)
	}
	return &user, nil
}

func (r *Repository) CreateUser(name string, password string, timezone string) (*models.User, error) {
	var user models.User
	err := psql.Insert("users").
		Columns("name", "password", "timezone").
		Values(name, password, timezone).
		Suffix("RETURNING name, password, timezone").
		RunWith(r.db).
		QueryRow().
		Scan(&user.Name, &user.Password, &user.Timezone)
	if err != nil {
		return nil, fmt.Errorf("create user: %v", err)
	}

	return &user, nil
}
