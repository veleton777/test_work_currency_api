package currency_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/veleton777/test_work_blum/internal/currency/v1"
	"github.com/veleton777/test_work_blum/internal/currency/v1/currency/entity"
	"github.com/veleton777/test_work_blum/internal/currency/v1/mocks"
	"github.com/veleton777/test_work_blum/internal/dto"
	"testing"
)

type CurrencyServiceTestSuite struct {
	suite.Suite
	svc               *currency.Svc
	mockCourseStorage *mocks.CourseStorage
	mockCurrencyRepo  *mocks.Repo
	mockCurrencyAPI   *mocks.CurrenciesAPI

	buf *bytes.Buffer
}

func (s *CurrencyServiceTestSuite) SetupTest() {
	s.buf = &bytes.Buffer{}
	l := zerolog.New(s.buf)

	s.mockCourseStorage = mocks.NewCourseStorage(s.T())
	s.mockCurrencyRepo = mocks.NewRepo(s.T())
	s.mockCurrencyAPI = mocks.NewCurrenciesAPI(s.T())
	s.svc = currency.NewCurrencySvc(
		s.mockCurrencyRepo,
		s.mockCurrencyAPI,
		s.mockCourseStorage,
		&l,
	)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyServiceTestSuite))
}

func (s *CurrencyServiceTestSuite) TestCreateCurrency_NoErr() {
	ctx := context.Background()

	s.mockCurrencyRepo.On("CreateCurrency", ctx, mock.Anything).
		Return(nil).Once()

	err := s.svc.CreateCurrency(ctx, dto.Currency{
		Name:        "test-1",
		Code:        "code-1",
		Type:        1,
		IsAvailable: true,
	})
	require.NoError(s.T(), err)

	s.mockCurrencyRepo.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestCreateCurrency_Err() {
	ctx := context.Background()

	testCases := []struct {
		name     string
		mockFunc func()
		req      dto.Currency
		expErr   error
	}{
		{
			name:     "invalid_currency_type",
			mockFunc: func() {},
			req: dto.Currency{
				Name:        "test-1",
				Code:        "code-1",
				Type:        555,
				IsAvailable: true,
			},
			expErr: errors.New("invalid currency type"),
		},
		{
			name: "currency_storage_err",
			mockFunc: func() {
				s.mockCurrencyRepo.On("CreateCurrency", ctx, mock.Anything).
					Return(errors.New("pg err")).Once()
			},
			req: dto.Currency{
				Name:        "test-1",
				Code:        "code-1",
				Type:        1,
				IsAvailable: true,
			},
			expErr: errors.New("pg err"),
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := s.svc.CreateCurrency(ctx, tc.req)
			require.Error(s.T(), err)
			require.ErrorContains(t, err, tc.expErr.Error())

			s.mockCurrencyRepo.AssertExpectations(s.T())
		})
	}
}

func (s *CurrencyServiceTestSuite) TestUpdateCurrency_NoErr() {
	ctx := context.Background()

	id := uuid.New()

	s.mockCurrencyRepo.On("UpdateCurrency", ctx, entity.Currency{
		ID:          id,
		Name:        "test-1",
		Code:        "code-1",
		Type:        1,
		IsAvailable: true,
	}).
		Return(nil).Once()

	err := s.svc.UpdateCurrency(ctx, dto.Currency{
		ID:          id,
		Name:        "test-1",
		Code:        "code-1",
		Type:        1,
		IsAvailable: true,
	})
	require.NoError(s.T(), err)

	s.mockCurrencyRepo.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestUpdateCurrency_Err() {
	ctx := context.Background()

	testID := uuid.New()

	testCases := []struct {
		name     string
		mockFunc func()
		req      dto.Currency
		expErr   error
	}{
		{
			name:     "invalid_currency_type",
			mockFunc: func() {},
			req: dto.Currency{
				Name:        "test-1",
				Code:        "code-1",
				Type:        555,
				IsAvailable: true,
			},
			expErr: errors.New("invalid currency type"),
		},
		{
			name: "currency_storage_err",
			mockFunc: func() {
				s.mockCurrencyRepo.On("UpdateCurrency", ctx, entity.Currency{
					ID:          testID,
					Name:        "test-1",
					Code:        "code-1",
					Type:        1,
					IsAvailable: true,
				}).
					Return(errors.New("pg err")).Once()
			},
			req: dto.Currency{
				ID:          testID,
				Name:        "test-1",
				Code:        "code-1",
				Type:        1,
				IsAvailable: true,
			},
			expErr: errors.New("pg err"),
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := s.svc.UpdateCurrency(ctx, tc.req)
			require.Error(s.T(), err)
			require.ErrorContains(t, err, tc.expErr.Error())

			s.mockCurrencyRepo.AssertExpectations(s.T())
		})
	}
}

func (s *CurrencyServiceTestSuite) TestDeleteCurrency_NoErr() {
	ctx := context.Background()

	id := uuid.New()

	s.mockCurrencyRepo.On("DeleteCurrency", ctx, id).
		Return(nil).Once()

	err := s.svc.DeleteCurrency(ctx, id)
	require.NoError(s.T(), err)

	s.mockCurrencyRepo.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestDeleteCurrency_Err() {
	ctx := context.Background()

	id := uuid.New()

	s.mockCurrencyRepo.On("DeleteCurrency", ctx, id).
		Return(errors.New("pg err")).Once()

	err := s.svc.DeleteCurrency(ctx, id)
	require.Error(s.T(), err)
	require.ErrorContains(s.T(), err, "pg err")

	s.mockCurrencyRepo.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestUpdateCourses_NoErr() {
	ctx := context.Background()

	s.mockCurrencyRepo.On("GetCurrencies", ctx).
		Return(entity.Currencies{
			{
				ID:          uuid.New(),
				Name:        "USD",
				Code:        "USD",
				Type:        2,
				IsAvailable: true,
			},
			{
				ID:          uuid.New(),
				Name:        "BTC",
				Code:        "BTC",
				Type:        1,
				IsAvailable: true,
			},
			{
				ID:          uuid.New(),
				Name:        "ETH",
				Code:        "ETH",
				Type:        1,
				IsAvailable: false,
			},
		}, nil).Once()

	s.mockCurrencyAPI.On("Convert", ctx, "USD", "BTC", float64(1)).
		Return(0.00045634, nil).Once()

	s.mockCurrencyAPI.On("Convert", ctx, "BTC", "USD", float64(1)).
		Return(float64(70000), nil).Once()

	s.mockCourseStorage.On("Set", ctx, "USD", "BTC", dto.CurrencyStorageDTO{
		Course:      decimal.NewFromFloat(0.00045634),
		IsAvailable: true,
	}).Return().Once()

	s.mockCourseStorage.On("Set", ctx, "BTC", "USD", dto.CurrencyStorageDTO{
		Course:      decimal.NewFromFloat(70000),
		IsAvailable: true,
	}).Return().Once()

	s.mockCourseStorage.On("Set", ctx, "ETH", "USD", dto.CurrencyStorageDTO{
		Course:      decimal.NewFromFloat(0),
		IsAvailable: false,
	}).Return().Once()

	s.mockCourseStorage.On("Set", ctx, "USD", "ETH", dto.CurrencyStorageDTO{
		Course:      decimal.NewFromFloat(0),
		IsAvailable: false,
	}).Return().Once()

	err := s.svc.UpdateCourses(ctx)
	require.NoError(s.T(), err)

	s.mockCourseStorage.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestUpdateCourses_StorageErr() {
	ctx := context.Background()

	s.mockCurrencyRepo.On("GetCurrencies", ctx).
		Return(nil, errors.New("pg err")).Once()

	err := s.svc.UpdateCourses(ctx)
	require.Error(s.T(), err)
	require.ErrorContains(s.T(), err, "pg err")

	s.mockCourseStorage.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestUpdateCourses_CurrencyAPIErr() {
	ctx := context.Background()

	s.mockCurrencyRepo.On("GetCurrencies", ctx).
		Return(entity.Currencies{
			{
				ID:          uuid.New(),
				Name:        "USD",
				Code:        "USD",
				Type:        2,
				IsAvailable: true,
			},
			{
				ID:          uuid.New(),
				Name:        "BTC",
				Code:        "BTC",
				Type:        1,
				IsAvailable: true,
			},
		}, nil).Once()

	s.mockCurrencyAPI.On("Convert", ctx, "USD", "BTC", float64(1)).
		Return(0.00045634, nil).Once()

	s.mockCurrencyAPI.On("Convert", ctx, "BTC", "USD", float64(1)).
		Return(float64(0), errors.New("api err")).Once()

	s.mockCourseStorage.On("Set", ctx, "USD", "BTC", dto.CurrencyStorageDTO{
		Course:      decimal.NewFromFloat(0.00045634),
		IsAvailable: true,
	}).Return().Once()

	s.mockCourseStorage.On("Set", ctx, "BTC", "USD", dto.CurrencyStorageDTO{
		Course:      decimal.NewFromFloat(0),
		IsAvailable: false,
	}).Return().Once()

	err := s.svc.UpdateCourses(ctx)
	require.NoError(s.T(), err)

	require.Contains(s.T(), string(s.buf.Bytes()), "convert currencies through api: from BTC to USD")
	require.Contains(s.T(), string(s.buf.Bytes()), "api err")

	s.mockCourseStorage.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestConvert_NoErr() {
	ctx := context.Background()
	data := dto.ConvertCurrencyReq{
		From:   "USD",
		To:     "BTC",
		Amount: 70000,
	}

	exp := decimal.NewFromFloat(0.00001441066)

	s.mockCourseStorage.On("Get", ctx, data.From, data.To).
		Return(exp, true).Once()

	res, err := s.svc.Convert(ctx, data)

	require.NoError(s.T(), err)
	require.Equal(s.T(), res, 1.0087462)

	s.mockCourseStorage.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestConvert_Err() {
	ctx := context.Background()
	data := dto.ConvertCurrencyReq{
		From:   "USD",
		To:     "BTC",
		Amount: 70000,
	}

	s.mockCourseStorage.On("Get", ctx, data.From, data.To).
		Return(decimal.Decimal{}, false).Once()

	res, err := s.svc.Convert(ctx, data)

	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, entity.ErrCurrencyNotAvailable)
	require.Equal(s.T(), res, float64(0))

	s.mockCourseStorage.AssertExpectations(s.T())
}
