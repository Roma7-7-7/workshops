package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Item is used by pop to map your items database table to your go code.
type Item struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Alias      string    `json:"alias" db:"alias"`
	Title      string    `json:"title" db:"title"`
	Desc       string    `json:"desc" db:"desc"`
	Pictures   string    `json:"pictures" db:"pictures"`
	Price      int       `json:"price" db:"price"`
	Count      int       `json:"count" db:"count"`
	CategoryID uuid.UUID `json:"category_id" db:"category_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (i Item) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Items is not required by pop and may be deleted
type Items []Item

// String is not required by pop and may be deleted
func (i Items) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (i *Item) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: i.Alias, Name: "Alias"},
		&validators.StringIsPresent{Field: i.Title, Name: "Title"},
		&validators.StringIsPresent{Field: i.Desc, Name: "Desc"},
		&validators.StringIsPresent{Field: i.Pictures, Name: "Pictures"},
		&validators.IntIsPresent{Field: i.Price, Name: "Price"},
		&validators.IntIsPresent{Field: i.Count, Name: "Count"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (i *Item) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (i *Item) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
