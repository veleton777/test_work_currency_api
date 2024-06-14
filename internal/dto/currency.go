package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Currency struct {
	ID   uuid.UUID `json:"-"`
	Name string    `json:"name" validate:"required" example:"Bitcoin"`
	Code string    `json:"code" validate:"required" example:"BTC"`
	// Type
	// * 1 - Crypto type
	// * 2 - Fiat type
	Type        int  `json:"type" validate:"required" enums:"1,2"`
	IsAvailable bool `json:"isAvailable" validate:"required"`
}

type CurrencyStorageDTO struct {
	Course      decimal.Decimal
	IsAvailable bool
}
