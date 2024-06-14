package main

import (
	"context"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/veleton777/test_work_blum/internal/config"
	"github.com/veleton777/test_work_blum/internal/server"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	l := zerolog.New(os.Stdout).
		Level(conf.LogLevel()).
		With().Timestamp().Stack().Caller().
		Logger()

	ctx := l.WithContext(context.Background())
	exitCode := 0

	s, err := server.New(ctx, &conf, &l)
	if err != nil {
		l.Fatal().Err(err).Msg("create server")

		exitCode = 1
		os.Exit(exitCode)
	}

	if err = s.Run(ctx); err != nil {
		err = errors.Wrap(err, "run server")
		l.Err(err).Send()

		exitCode = 1
	}

	os.Exit(exitCode)
}
