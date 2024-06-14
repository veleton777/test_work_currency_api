//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"github.com/veleton777/test_work_blum/internal/config"
	"github.com/veleton777/test_work_blum/internal/currency/v1/currency/entity"
	"github.com/veleton777/test_work_blum/internal/currency/v1/currency/storage/postgres"
	"github.com/veleton777/test_work_blum/internal/currency/v1/currency/storage/postgres/converter"
	storageentity "github.com/veleton777/test_work_blum/internal/currency/v1/currency/storage/postgres/entity"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	repo      *postgres.RepoPostgres
	pgxClient *pgxpool.Pool
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	ctx := context.Background()
	conf, err := config.Load()
	s.Require().NoError(err)

	pgCfg, err := pgxpool.ParseConfig(
		fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s",
			conf.PgHost(),
			conf.PgPort(),
			conf.PgDB(),
			conf.PgUser(),
			conf.PgPassword(),
		),
	)
	s.Require().NoError(err)

	pgClient, err := pgxpool.NewWithConfig(ctx, pgCfg)
	s.Require().NoError(err)

	err = pgClient.Ping(ctx)
	s.Require().NoError(err)

	s.repo = postgres.NewRepoPostgres(pgClient, 1*time.Second)
	s.pgxClient = pgClient

	s.clearCollection()
}

func (s *Suite) TearDownSuite() {
	s.clearCollection()
}

func (s *Suite) afterEachTest() {
	s.clearCollection()
}

func (s *Suite) clearCollection() {
	ctx := context.Background()

	_, err := s.pgxClient.Exec(ctx, "TRUNCATE TABLE currencies")
	s.Require().NoError(err)
}

func (s *Suite) TearDownTest() {
	s.afterEachTest()
}

func (s *Suite) TearDownSubTest() {
	s.afterEachTest()
}

func (s *Suite) TestGetCurrencies_NoErr() {
	ctx := context.Background()

	currency, err := s.createCurrency(ctx)
	s.Require().NoError(err)

	currencies, err := s.repo.GetCurrencies(ctx)
	s.Require().NoError(err)

	s.Require().Equal(len(currencies), 1)

	s.Require().Equal(currency.ID, currencies[0].ID)
	s.Require().Equal(currency.Name, currencies[0].Name)
	s.Require().Equal(currency.Code, currencies[0].Code)
	s.Require().Equal(currency.Type, currencies[0].Type)
	s.Require().Equal(currency.IsAvailable, currencies[0].IsAvailable)
}

func (s *Suite) TestCreateCurrency_NoErr() {
	ctx := context.Background()

	currencies, err := s.currencies(ctx)
	s.Require().NoError(err)

	s.Require().Equal(len(currencies), 0)

	id := uuid.New()

	err = s.repo.CreateCurrency(ctx, entity.Currency{
		ID:          id,
		Name:        "name-1",
		Code:        "code-1",
		Type:        1,
		IsAvailable: true,
	})
	s.Require().NoError(err)

	currencies, err = s.currencies(ctx)
	s.Require().NoError(err)

	s.Require().Equal(len(currencies), 1)

	s.Require().Equal(id, currencies[0].ID)
	s.Require().Equal("name-1", currencies[0].Name)
	s.Require().Equal("code-1", currencies[0].Code)
	s.Require().Equal(1, currencies[0].Type)
	s.Require().Equal(true, currencies[0].IsAvailable)
}

func (s *Suite) TestUpdateCurrency_NoErr() {
	ctx := context.Background()

	currency, err := s.createCurrency(ctx)
	s.Require().NoError(err)

	currency.Name = "name-2"
	currency.Code = "code-2"
	currency.Type = 2
	currency.IsAvailable = false

	err = s.repo.UpdateCurrency(ctx, currency)
	s.Require().NoError(err)

	currencies, err := s.currencies(ctx)
	s.Require().NoError(err)

	s.Require().Equal(len(currencies), 1)

	s.Require().Equal(currency.ID, currencies[0].ID)
	s.Require().Equal("name-2", currencies[0].Name)
	s.Require().Equal("code-2", currencies[0].Code)
	s.Require().Equal(2, currencies[0].Type)
	s.Require().Equal(false, currencies[0].IsAvailable)
}

func (s *Suite) TestUpdateCurrency_ReturnNotFoundErr() {
	ctx := context.Background()

	err := s.repo.UpdateCurrency(ctx, entity.Currency{
		ID:          uuid.New(),
		Name:        "test-1",
		Code:        "code-1",
		Type:        1,
		IsAvailable: true,
	})
	s.Require().Error(err)
	s.Require().ErrorIs(err, entity.ErrEntityNotFound)
}

func (s *Suite) TestDeleteCurrency_NoErr() {
	ctx := context.Background()

	currency, err := s.createCurrency(ctx)
	s.Require().NoError(err)

	err = s.repo.DeleteCurrency(ctx, currency.ID)
	s.Require().NoError(err)

	currencies, err := s.currencies(ctx)
	s.Require().NoError(err)

	s.Require().Equal(len(currencies), 0)
}

func (s *Suite) TestDeleteCurrency_ReturnNotFoundErr() {
	ctx := context.Background()

	err := s.repo.DeleteCurrency(ctx, uuid.New())
	s.Require().Error(err)
	s.Require().ErrorIs(err, entity.ErrEntityNotFound)
}

func (s *Suite) currencies(ctx context.Context) (storageentity.Currencies, error) {
	rows, err := s.pgxClient.Query(ctx, "SELECT id, name, code, type, is_available FROM currencies")
	if err != nil {
		return nil, errors.Wrap(err, "pgx query")
	}

	defer rows.Close()

	currencies, err := pgx.CollectRows(rows, pgx.RowToStructByName[storageentity.Currency])
	if err != nil {
		return nil, errors.Wrap(err, "currencies to structs")
	}

	return currencies, nil
}

func (s *Suite) createCurrency(ctx context.Context) (entity.Currency, error) {
	currencies, err := s.currencies(ctx)
	if err != nil {
		return entity.Currency{}, errors.Wrap(err, "get currencies")
	}

	if len(currencies) != 0 {
		return entity.Currency{}, errors.New("quantity currencies > 0")
	}

	id := uuid.New()

	err = s.repo.CreateCurrency(ctx, entity.Currency{
		ID:          id,
		Name:        "name-1",
		Code:        "code-1",
		Type:        1,
		IsAvailable: true,
	})
	if err != nil {
		return entity.Currency{}, errors.Wrap(err, "create currency in postgres")
	}

	currencies, err = s.currencies(ctx)
	if err != nil {
		return entity.Currency{}, errors.Wrap(err, "get currencies")
	}

	if len(currencies) != 1 {
		return entity.Currency{}, errors.New("quantity currencies != 1")
	}

	res, err := converter.CurrencyToEntity(currencies[0])
	if err != nil {
		return entity.Currency{}, errors.Wrap(err, "convert currency to entity")
	}

	return res, nil
}
