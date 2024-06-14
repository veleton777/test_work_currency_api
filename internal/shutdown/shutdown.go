package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
)

type (
	Shutdown   [3][]shutdownFn
	shutdownFn func(context.Context) error
)

func New() *Shutdown {
	return &Shutdown{}
}

func (s *Shutdown) WaitShutdown(ctx context.Context) {
	stopSignals := []os.Signal{syscall.SIGTERM, syscall.SIGINT}
	sig := make(chan os.Signal, len(stopSignals))
	signal.Notify(sig, stopSignals...)
	zerolog.Ctx(ctx).Info().Msgf("got %s os signal. application will be stopped", <-sig)

	s.do(ctx)
}

func (s *Shutdown) AddHiPriority(fn shutdownFn) {
	s[0] = append(s[0], fn)
}

func (s *Shutdown) AddNormalPriority(fn shutdownFn) {
	s[1] = append(s[1], fn)
}

func (s *Shutdown) AddLowPriority(fn shutdownFn) {
	s[2] = append(s[2], fn)
}

func (s *Shutdown) do(ctx context.Context) {
	for _, priorityShutdown := range s {
		for _, fn := range priorityShutdown {
			if err := fn(ctx); err != nil {
				zerolog.Ctx(ctx).Err(err).Send()
			}
		}
	}
}
