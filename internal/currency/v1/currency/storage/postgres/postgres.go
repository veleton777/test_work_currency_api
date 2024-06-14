package postgres

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/veleton777/test_work_blum/internal/currency/v1/currency/entity"
	"github.com/veleton777/test_work_blum/internal/currency/v1/currency/storage/postgres/converter"
	storageentity "github.com/veleton777/test_work_blum/internal/currency/v1/currency/storage/postgres/entity"
)

const (
	currenciesTable = "currencies"

	pgxDuplicateKeyCode = "23505"
)

type RepoPostgres struct {
	pgClient *pgxpool.Pool
	timeout  time.Duration
}

func NewRepoPostgres(pgClient *pgxpool.Pool, timeout time.Duration) *RepoPostgres {
	return &RepoPostgres{pgClient: pgClient, timeout: timeout}
}

func (r *RepoPostgres) GetCurrencies(ctx context.Context) (entity.Currencies, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := squirrel.Select("id", "name", "code", "type", "is_available").
		From(currenciesTable)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query to sql")
	}

	rows, err := r.pgClient.Query(ctx, query, v...)
	if err != nil {
		return nil, errors.Wrap(err, "pgx query")
	}
	defer rows.Close()

	currencies, err := pgx.CollectRows(rows, pgx.RowToStructByName[storageentity.Currency])
	if err != nil {
		return nil, errors.Wrap(err, "scan resp to struct")
	}

	res, err := converter.CurrenciesToEntity(currencies)
	if err != nil {
		return nil, errors.Wrap(err, "convert currencies to entity")
	}

	return res, nil
}

func (r *RepoPostgres) GetCurrencyByCode(ctx context.Context, code string) (entity.Currency, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := squirrel.Select("id", "name", "code", "type", "is_available").
		From(currenciesTable).
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"code": code})

	query, v, err := builder.ToSql()
	if err != nil {
		return entity.Currency{}, errors.Wrap(err, "query to sql")
	}

	rows, err := r.pgClient.Query(ctx, query, v...)
	if err != nil {
		return entity.Currency{}, errors.Wrap(err, "pgx query")
	}

	currency, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[storageentity.Currency])
	if err != nil {
		return entity.Currency{}, errors.Wrap(err, "scan resp to struct")
	}

	res, err := converter.CurrencyToEntity(currency)
	if err != nil {
		return entity.Currency{}, errors.Wrap(err, "convert currency to entity")
	}

	return res, nil
}

func (r *RepoPostgres) CreateCurrency(ctx context.Context, currency entity.Currency) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := squirrel.Insert(currenciesTable).
		PlaceholderFormat(squirrel.Dollar).
		Columns("id", "name", "code", "type", "is_available").
		Values(currency.ID, currency.Name, currency.Code, currency.Type, currency.IsAvailable)

	query, v, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "query to sql")
	}

	if _, err = r.pgClient.Exec(ctx, query, v...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgxDuplicateKeyCode {
				return entity.ErrCurrencyAlreadyExists
			}
		}

		return errors.Wrap(err, "exec pg query")
	}

	return nil
}

func (r *RepoPostgres) UpdateCurrency(ctx context.Context, currency entity.Currency) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := squirrel.Update(currenciesTable).
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"id": currency.ID}).
		Set("name", currency.Name).
		Set("code", currency.Code).
		Set("type", currency.Type).
		Set("is_available", currency.IsAvailable)

	query, v, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "query to sql")
	}

	cmd, err := r.pgClient.Exec(ctx, query, v...)
	if err != nil {
		return errors.Wrap(err, "exec pg query")
	}

	if cmd.RowsAffected() == 0 {
		return entity.ErrEntityNotFound
	}

	return nil
}

func (r *RepoPostgres) DeleteCurrency(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := squirrel.Delete(currenciesTable).
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"id": id})

	query, v, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "query to sql")
	}

	cmd, err := r.pgClient.Exec(ctx, query, v...)
	if err != nil {
		return errors.Wrap(err, "exec pg query")
	}

	if cmd.RowsAffected() == 0 {
		return entity.ErrEntityNotFound
	}

	return nil
}
