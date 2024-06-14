package currency

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/veleton777/test_work_blum/internal/currency/v1/currency/entity"
	"github.com/veleton777/test_work_blum/internal/dto"
)

//go:generate mockery --name Repo
type Repo interface {
	GetCurrencies(ctx context.Context) (entity.Currencies, error)
	CreateCurrency(ctx context.Context, currency entity.Currency) error
	UpdateCurrency(ctx context.Context, currency entity.Currency) error
	DeleteCurrency(ctx context.Context, id uuid.UUID) error
}

//go:generate mockery --name CurrenciesAPI
type CurrenciesAPI interface {
	Convert(ctx context.Context, from, to string, amount float64) (float64, error)
}

//go:generate mockery --name CourseStorage
type CourseStorage interface {
	Set(ctx context.Context, codeFrom, codeTo string, data dto.CurrencyStorageDTO)
	Get(ctx context.Context, codeFrom, codeTo string) (decimal.Decimal, bool)
}

//go:generate mockery --name Storage
type Storage interface {
	Set(codeFrom, codeTo string, data dto.CurrencyStorageDTO)
	Get(codeFrom, codeTo string) (decimal.Decimal, bool)
}
