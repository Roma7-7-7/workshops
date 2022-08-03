package postgre

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"time"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	db *sql.DB
}

func (r *Repository) GetEvents(title, dateFrom, timeFrom, dateTo, timeTo string) ([]*models.Event, error) {
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
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf(`querying with sql="%s": %v`, query, err)
	}

	var result []*models.Event
	for rows.Next() {
		var event models.Event
		if err = rows.Scan(&event.ID, &event.Title, &event.Description, &event.TimeFrom, &event.TimeTo, pq.Array(&event.Notes)); err != nil {
			return nil, fmt.Errorf("scan event: %v", err)
		}
		result = append(result, &event)
	}

	return result, nil
}

func (r *Repository) GetEvent(id string) (*models.Event, error) {
	query, args, err := psql.Select("id", "title", "description", "timestamp_from", "timestamp_to", "notes").
		From("event").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build event query: %v", err)
	}
	var event models.Event
	if err = r.db.QueryRow(query, args...).Scan(&event.ID, &event.Title, &event.Description, &event.TimeFrom, &event.TimeTo, pq.Array(&event.Notes)); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("scan event: %v", err)
	}

	return &event, nil
}

func (r *Repository) Exists(id string) (bool, error) {
	var count int
	if err := psql.Select("COUNT(*)").From("event").Where(sq.Eq{"id": id}).RunWith(r.db).QueryRow().Scan(&count); err != nil {
		return false, fmt.Errorf("check if exists: %v", err)
	}
	return count > 0, nil
}

func (r *Repository) CreateEvent(title string, description string, from time.Time, to time.Time, notes []string) (*models.Event, error) {
	query := psql.Insert("event").
		Columns("id", "title", "description", "timestamp_from", "timestamp_to", "notes").
		Values(sq.Expr("gen_random_uuid()"), title, description, from, to, pq.Array(notes)).
		Suffix("RETURNING id, title, description, timestamp_from, timestamp_to, notes").
		RunWith(r.db)

	var event models.Event

	err := query.QueryRow().Scan(&event.ID, &event.Title, &event.Description, &event.TimeFrom, &event.TimeTo, pq.Array(&event.Notes))
	if err != nil {
		return nil, fmt.Errorf("insert event: %v", err)
	}

	return &event, nil
}

func (r *Repository) UpdateEvent(id, title, description string, from time.Time, to time.Time, notes []string) (*models.Event, error) {
	query := psql.Update("event").
		Set("title", title).
		Set("description", description).
		Set("timestamp_from", from).
		Set("timestamp_to", to).
		Set("notes", pq.Array(notes)).
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		Suffix("RETURNING id, title, description, timestamp_from, timestamp_to, notes")

	var event models.Event

	err := query.QueryRow().Scan(&event.ID, &event.Title, &event.Description, &event.TimeFrom, &event.TimeTo, pq.Array(&event.Notes))
	if err != nil {
		return nil, fmt.Errorf("update event: %v", err)
	}

	return &event, nil
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
