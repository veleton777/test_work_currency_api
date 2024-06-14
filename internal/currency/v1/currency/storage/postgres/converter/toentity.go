package converter

import (
	"github.com/pkg/errors"
	"github.com/veleton777/test_work_blum/internal/currency/v1/currency/entity"
	storageentity "github.com/veleton777/test_work_blum/internal/currency/v1/currency/storage/postgres/entity"
)

func CurrencyToEntity(cur storageentity.Currency) (entity.Currency, error) {
	t, err := entity.IntToCurrencyType(cur.Type)
	if err != nil {
		return entity.Currency{}, errors.Wrap(err, "int to currency type")
	}

	return entity.Currency{
		ID:          cur.ID,
		Name:        cur.Name,
		Code:        cur.Code,
		Type:        t,
		IsAvailable: cur.IsAvailable,
	}, nil
}

func CurrenciesToEntity(cur []storageentity.Currency) (entity.Currencies, error) {
	res := make([]entity.Currency, 0, len(cur))

	for _, c := range cur {
		e, err := CurrencyToEntity(c)
		if err != nil {
			return nil, errors.Wrap(err, "convert currency to entity")
		}

		res = append(res, e)
	}

	return res, nil
}
