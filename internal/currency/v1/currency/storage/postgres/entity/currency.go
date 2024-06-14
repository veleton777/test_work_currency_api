package entity

import "github.com/google/uuid"

type Currency struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Code        string    `db:"code"`
	Type        int       `db:"type"`
	IsAvailable bool      `db:"is_available"`
}

type Currencies []Currency
