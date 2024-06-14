package entity

import "errors"

var (
	ErrEntityNotFound        = errors.New("entity not found")
	ErrInvalidCurrencyType   = errors.New("invalid currency type")
	ErrCurrencyNotAvailable  = errors.New("currency not available")
	ErrCurrencyAlreadyExists = errors.New("currency already exists")
)
