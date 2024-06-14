package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/shopspring/decimal"
	"github.com/veleton777/test_work_blum/internal/dto"
)

type value struct {
	course      decimal.Decimal
	isAvailable bool
}

type Storage struct {
	storage map[string]value
	mu      *sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		storage: make(map[string]value),
		mu:      &sync.RWMutex{},
	}
}

func (s *Storage) Set(_ context.Context, codeFrom, codeTo string, data dto.CurrencyStorageDTO) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage[s.key(codeFrom, codeTo)] = value{
		course:      data.Course,
		isAvailable: data.IsAvailable,
	}
}

func (s *Storage) Get(_ context.Context, codeFrom, codeTo string) (decimal.Decimal, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.storage[s.key(codeFrom, codeTo)]
	if !ok {
		return decimal.Decimal{}, false
	}

	if !v.isAvailable {
		return decimal.Decimal{}, false
	}

	return v.course, true
}

func (s *Storage) key(codeFrom, codeTo string) string {
	return fmt.Sprintf("%s_%s", codeFrom, codeTo)
}
