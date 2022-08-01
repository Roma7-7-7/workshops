package postgre

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	db *sql.DB
}

func (r Repository) GetEvents(title, dateFrom, timeFrom, dateTo, timeTo string) ([]models.Event, error) {
	filters := sq.And{}

	if title != "" {
		filters = append(filters, sq.Eq{"title": title})
	}
	if dateFrom != "" {
		filters = append(filters, sq.GtOrEq{"timestamp_from::date": dateFrom})
	}
	if timeFrom != "" {
		filters = append(filters, sq.GtOrEq{"timestamp_from::time": timeFrom})
	}
	if dateTo != "" {
		filters = append(filters, sq.GtOrEq{"timestamp_to::date": dateTo})
	}
	if timeTo != "" {
		filters = append(filters, sq.GtOrEq{"timestamp_to::time": timeTo})
	}

	q := psql.Select("id", "title", "description", "timestamp_from", "timestamp_to", "notes").From("event")
	if len(filters) > 0 {
		q = q.Where(filters)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build events query: %v", err)
	}

	var rows *sql.Rows
	rows, err = r.db.Query(query, args...)

	if err != nil {
		return nil, fmt.Errorf(`querying with sql="%s": %v`, query, err)
	}

	var result []models.Event
	for rows.Next() {
		var event models.Event
		if err = rows.Scan(&event.ID, &event.Title, &event.Description, &event.TimeFrom, &event.TimeTo, pq.Array(&event.Notes)); err != nil {
			return nil, fmt.Errorf("scan event: %v", err)
		}
		result = append(result, event)
	}

	return result, nil
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
