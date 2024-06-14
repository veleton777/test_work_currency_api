package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/veleton777/test_work_blum/internal/common"
	"github.com/veleton777/test_work_blum/internal/config"
	"github.com/veleton777/test_work_blum/internal/currency/v1"
	"github.com/veleton777/test_work_blum/internal/currency/v1/currency/storage/memory"
	"github.com/veleton777/test_work_blum/internal/currency/v1/currency/storage/postgres"
	"github.com/veleton777/test_work_blum/internal/pkg/fastforex"
	"github.com/veleton777/test_work_blum/internal/shutdown"
	v1 "github.com/veleton777/test_work_blum/internal/transport/http/v1"
)

type API struct {
	config *config.Config
	sh     *shutdown.Shutdown
	l      *zerolog.Logger

	currencyServer *v1.CurrencyServer
	currencySvc    *currency.Svc
}

func New(ctx context.Context, config *config.Config, l *zerolog.Logger) (*API, error) {
	sh := shutdown.New()

	a := &API{ //nolint:exhaustruct
		config: config,
		sh:     sh,
		l:      l,
	}

	pgxClient, err := a.pgxClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "create pgx client")
	}

	currencyRepo := postgres.NewRepoPostgres(pgxClient, a.config.PgTimeout())
	courseStorage := memory.NewStorage()

	httpClient := &http.Client{ //nolint:exhaustruct
		Timeout: a.config.FastForexHTTPTimeout(),
	}

	fastForexClient := fastforex.NewClient(
		a.config.FastForexAPIHost(),
		a.config.FastForexAPIKey(),
		httpClient,
	)

	currencySvc := currency.NewCurrencySvc(currencyRepo, fastForexClient, courseStorage, l)
	a.currencySvc = currencySvc

	a.currencyServer = v1.NewCurrencyServer(currencySvc)

	return a, nil
}

func (s *API) Run(ctx context.Context) error {
	var err error

	app := fiber.New(fiber.Config{ //nolint:exhaustruct
		DisableStartupMessage: true,
	})

	app.Use(fiberzerolog.New(fiberzerolog.Config{ //nolint:exhaustruct
		Logger: s.l,
	}))

	s.routes(app)

	s.sh.AddNormalPriority(func(_ context.Context) error {
		if err := app.Shutdown(); err != nil {
			return errors.Wrap(err, "fiber app shutdown")
		}

		return nil
	})

	if err = s.currencySvc.UpdateCourses(ctx); err != nil {
		return errors.Wrap(err, "first update courses")
	}

	go common.BackgroundWorker(ctx, s.config.FastForexBackgroundTaskDelay(), s.l, s.currencySvc.UpdateCourses)

	go func() {
		s.l.Info().Msg("start server")

		if err = app.Listen(s.config.HTTPAddr()); err != nil {
			err = errors.Wrap(err, "fiber app listen port")
		}
	}()

	s.sh.WaitShutdown(ctx)

	return err
}

func (s *API) pgxClient(ctx context.Context) (*pgxpool.Pool, error) {
	pgCfg, err := pgxpool.ParseConfig(
		fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s",
			s.config.PgHost(),
			s.config.PgPort(),
			s.config.PgDB(),
			s.config.PgUser(),
			s.config.PgPassword(),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "create pgxpool config")
	}

	pgClient, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, errors.Wrap(err, "connect to pg")
	}

	if err = pgClient.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "pg ping")
	}

	s.sh.AddNormalPriority(func(_ context.Context) error {
		pgClient.Close()

		return nil
	})

	return pgClient, nil
}
