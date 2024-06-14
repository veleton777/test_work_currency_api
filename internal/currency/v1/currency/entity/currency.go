package entity

import (
	"github.com/google/uuid"
)

type Currency struct {
	ID          uuid.UUID
	Name        string
	Code        string
	Type        CurrencyType
	IsAvailable bool
}

type Currencies []Currency

type CurrencyType int

const (
	TypeCrypto CurrencyType = 1
	TypeFiat   CurrencyType = 2
)

func IntToCurrencyType(t int) (CurrencyType, error) {
	switch CurrencyType(t) {
	case TypeCrypto:
		return TypeCrypto, nil
	case TypeFiat:
		return TypeFiat, nil
	}

	return 0, ErrInvalidCurrencyType
}

func (c Currency) IsFiat() bool {
	return c.Type == TypeFiat
}

func (c Currency) IsCrypto() bool {
	return c.Type == TypeCrypto
}
