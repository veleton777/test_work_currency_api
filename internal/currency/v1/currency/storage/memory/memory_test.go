//nolint:testpackage
package memory

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/veleton777/test_work_blum/internal/dto"
	"testing"
)

type MemoryStorageTestSuite struct {
	suite.Suite
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryStorageTestSuite))
}

func (s *MemoryStorageTestSuite) TestSetMethod() {
	ctx := context.Background()
	st := NewStorage()

	require.Equal(s.T(), len(st.storage), 0)

	st.Set(ctx, "USD", "BTC", dto.CurrencyStorageDTO{
		Course:      decimal.NewFromFloat(123.567),
		IsAvailable: true,
	})

	require.Equal(s.T(), len(st.storage), 1)
}

func (s *MemoryStorageTestSuite) TestGetMethod() {
	ctx := context.Background()
	st := NewStorage()

	v, ok := st.Get(ctx, "USD", "BTC")

	require.Equal(s.T(), v, decimal.Decimal{})
	require.False(s.T(), ok)

	course := decimal.NewFromFloat(123.567)

	st.Set(ctx, "USD", "BTC", dto.CurrencyStorageDTO{
		Course:      course,
		IsAvailable: true,
	})

	v, ok = st.Get(ctx, "USD", "BTC")

	require.Equal(s.T(), v, course)
	require.True(s.T(), ok)
}

func (s *MemoryStorageTestSuite) TestGetMethod_IsNotAvailable() {
	ctx := context.Background()
	st := NewStorage()

	v, ok := st.Get(ctx, "USD", "BTC")

	require.Equal(s.T(), v, decimal.Decimal{})
	require.False(s.T(), ok)

	course := decimal.NewFromFloat(123.567)

	st.Set(ctx, "USD", "BTC", dto.CurrencyStorageDTO{
		Course:      course,
		IsAvailable: false,
	})

	v, ok = st.Get(ctx, "USD", "BTC")

	require.Equal(s.T(), v, decimal.Decimal{})
	require.False(s.T(), ok)
}
