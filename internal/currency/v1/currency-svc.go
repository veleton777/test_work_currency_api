package currency

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/veleton777/test_work_blum/internal/currency/v1/currency/entity"
	"github.com/veleton777/test_work_blum/internal/dto"
)

type Svc struct {
	currencyStorage Repo
	currenciesAPI   CurrenciesAPI
	courseStorage   CourseStorage
	wg              *sync.WaitGroup
	l               *zerolog.Logger
}

func NewCurrencySvc(
	currencyStorage Repo,
	currenciesAPI CurrenciesAPI,
	courseStorage CourseStorage,
	l *zerolog.Logger,
) *Svc {
	return &Svc{
		currencyStorage: currencyStorage,
		currenciesAPI:   currenciesAPI,
		courseStorage:   courseStorage,
		wg:              &sync.WaitGroup{},
		l:               l,
	}
}

func (s *Svc) CreateCurrency(ctx context.Context, dto dto.Currency) error {
	t, err := entity.IntToCurrencyType(dto.Type)
	if err != nil {
		return errors.Wrap(err, "int to currency type")
	}

	currency := entity.Currency{
		ID:          uuid.New(),
		Name:        dto.Name,
		Code:        dto.Code,
		Type:        t,
		IsAvailable: dto.IsAvailable,
	}

	if err = s.currencyStorage.CreateCurrency(ctx, currency); err != nil {
		return errors.Wrap(err, "save currency to storage")
	}

	return nil
}

func (s *Svc) UpdateCurrency(ctx context.Context, dto dto.Currency) error {
	t, err := entity.IntToCurrencyType(dto.Type)
	if err != nil {
		return errors.Wrap(err, "int to currency type")
	}

	currency := entity.Currency{
		ID:          dto.ID,
		Name:        dto.Name,
		Code:        dto.Code,
		Type:        t,
		IsAvailable: dto.IsAvailable,
	}

	if err = s.currencyStorage.UpdateCurrency(ctx, currency); err != nil {
		return errors.Wrap(err, "update currency in storage")
	}

	return nil
}

func (s *Svc) DeleteCurrency(ctx context.Context, id uuid.UUID) error {
	if err := s.currencyStorage.DeleteCurrency(ctx, id); err != nil {
		return errors.Wrap(err, "delete currency from storage")
	}

	return nil
}

func (s *Svc) Convert(ctx context.Context, dto dto.ConvertCurrencyReq) (float64, error) {
	v, ok := s.courseStorage.Get(ctx, dto.From, dto.To)
	if !ok {
		return 0, entity.ErrCurrencyNotAvailable
	}

	res, _ := v.Mul(decimal.NewFromFloat(dto.Amount)).Float64()

	return res, nil
}

func (s *Svc) UpdateCourses(ctx context.Context) error {
	var (
		fiatCur   entity.Currencies
		cryptoCur entity.Currencies
	)

	currencies, err := s.currencyStorage.GetCurrencies(ctx)
	if err != nil {
		return errors.Wrap(err, "get currencies from storage")
	}

	for _, c := range currencies {
		if c.IsFiat() {
			fiatCur = append(fiatCur, c)

			continue
		}

		cryptoCur = append(cryptoCur, c)
	}

	for _, f := range fiatCur {
		for _, c := range cryptoCur {
			s.wg.Add(2) //nolint:mnd
			go s.updateCourse(ctx, f, c, s.wg)
			go s.updateCourse(ctx, c, f, s.wg)
		}
	}

	s.wg.Wait()

	return nil
}

func (s *Svc) updateCourse(ctx context.Context, from entity.Currency, to entity.Currency, wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		course      float64
		isAvailable bool
		err         error
	)

	if from.IsAvailable && to.IsAvailable {
		isAvailable = true
	}

	if isAvailable {
		course, err = s.currenciesAPI.Convert(ctx, from.Code, to.Code, 1)
		if err != nil {
			isAvailable = false

			s.l.Err(err).Msgf("convert currencies through api: from %s to %s", from.Code, to.Code)
		}
	}

	s.courseStorage.Set(ctx, from.Code, to.Code, dto.CurrencyStorageDTO{
		Course:      decimal.NewFromFloat(course),
		IsAvailable: isAvailable,
	})
}
