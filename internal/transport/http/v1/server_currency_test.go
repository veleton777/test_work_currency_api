//go:build integration

package v1_test

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/veleton777/test_work_blum/internal/currency/v1/currency/entity"
	v1 "github.com/veleton777/test_work_blum/internal/transport/http/v1"
	"github.com/veleton777/test_work_blum/internal/transport/http/v1/mocks"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

type ServerCurrencySuite struct {
	suite.Suite

	srv             *v1.CurrencyServer
	mockCurrencySvc *mocks.CurrencySvc
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ServerCurrencySuite))
}

func (s *ServerCurrencySuite) SetupSuite() {
	s.mockCurrencySvc = mocks.NewCurrencySvc(s.T())

	s.srv = v1.NewCurrencyServer(s.mockCurrencySvc)
}

func (s *ServerCurrencySuite) TestCreateCurrency() {
	ctx := context.Background()

	testCases := []struct {
		name     string
		data     string
		mockFunc func()
		expRes   string
		expCode  int
	}{
		{
			name: "success",
			data: `{"name": "test", "code": "code", "type": 1, "isAvailable": true}`,
			mockFunc: func() {
				s.mockCurrencySvc.On("CreateCurrency", ctx, mock.Anything).
					Return(nil).Once()
			},
			expRes:  "",
			expCode: 201,
		},
		{
			name:     "invalid_json",
			data:     `invalid_json`,
			mockFunc: func() {},
			expRes:   `{"code":400,"text":"invalid json body format"}`,
			expCode:  400,
		},
		{
			name:     "validation_err",
			data:     `{"name": "test", "type": 1, "isAvailable": true}`,
			mockFunc: func() {},
			expRes:   `{"code":400,"text":"Key: 'Currency.Code' Error:Field validation for 'Code' failed on the 'required' tag"}`,
			expCode:  400,
		},
		{
			name: "currency_already_exists",
			data: `{"name": "test", "code": "code", "type": 1, "isAvailable": true}`,
			mockFunc: func() {
				s.mockCurrencySvc.On("CreateCurrency", ctx, mock.Anything).
					Return(entity.ErrCurrencyAlreadyExists).Once()
			},
			expRes:  `{"code":400,"text":"currency already exists"}`,
			expCode: 400,
		},
		{
			name: "svc_err",
			data: `{"name": "test", "code": "code", "type": 1, "isAvailable": true}`,
			mockFunc: func() {
				s.mockCurrencySvc.On("CreateCurrency", ctx, mock.Anything).
					Return(errors.New("")).Once()
			},
			expRes:  `{"code":500,"text":"Internal Server error"}`,
			expCode: 500,
		},
	}

	for _, c := range testCases {
		s.T().Run(c.name, func(t *testing.T) {
			app := fiber.New()
			app.Post("/", s.srv.CreateCurrency)

			c.mockFunc()

			req := httptest.NewRequest("POST", "/", strings.NewReader(c.data))

			resp, err := app.Test(req, 1)
			s.Require().NoError(err)

			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			s.Require().NoError(err)

			assert.Equal(s.T(), c.expCode, resp.StatusCode)
			assert.Contains(s.T(), string(respBody), c.expRes)
		})
	}
}

func (s *ServerCurrencySuite) TestUpdateCurrency() {
	ctx := context.Background()

	testCases := []struct {
		name     string
		id       string
		data     string
		mockFunc func()
		expRes   string
		expCode  int
	}{
		{
			name: "success",
			id:   uuid.New().String(),
			data: `{"name": "test", "code": "code", "type": 1, "isAvailable": true}`,
			mockFunc: func() {
				s.mockCurrencySvc.On("UpdateCurrency", ctx, mock.Anything).
					Return(nil).Once()
			},
			expRes:  "",
			expCode: 204,
		},
		{
			name:     "invalid_id",
			id:       "123",
			data:     `{"name": "test", "code": "code", "type": 1, "isAvailable": true}`,
			mockFunc: func() {},
			expRes:   `{"code":400,"text":"invalid id format"}`,
			expCode:  400,
		},
		{
			name:     "invalid_json",
			id:       uuid.New().String(),
			data:     `invalid_json`,
			mockFunc: func() {},
			expRes:   `{"code":400,"text":"invalid json body format"}`,
			expCode:  400,
		},
		{
			name:     "validation_err",
			id:       uuid.New().String(),
			data:     `{"name": "test", "type": 1, "isAvailable": true}`,
			mockFunc: func() {},
			expRes:   `{"code":400,"text":"Key: 'Currency.Code' Error:Field validation for 'Code' failed on the 'required' tag"}`,
			expCode:  400,
		},
		{
			name: "currency_not_found",
			id:   uuid.New().String(),
			data: `{"name": "test", "code": "code", "type": 1, "isAvailable": true}`,
			mockFunc: func() {
				s.mockCurrencySvc.On("UpdateCurrency", ctx, mock.Anything).
					Return(entity.ErrEntityNotFound).Once()
			},
			expRes:  `{"code":404,"text":"Not Found"}`,
			expCode: 404,
		},
		{
			name: "svc_err",
			id:   uuid.New().String(),
			data: `{"name": "test", "code": "code", "type": 1, "isAvailable": true}`,
			mockFunc: func() {
				s.mockCurrencySvc.On("UpdateCurrency", ctx, mock.Anything).
					Return(errors.New("")).Once()
			},
			expRes:  `{"code":500,"text":"Internal Server error"}`,
			expCode: 500,
		},
	}

	for _, c := range testCases {
		s.T().Run(c.name, func(t *testing.T) {
			app := fiber.New()
			app.Put("/:id", s.srv.UpdateCurrency)

			c.mockFunc()

			req := httptest.NewRequest("PUT", "/"+c.id, strings.NewReader(c.data))

			resp, err := app.Test(req, 1)
			s.Require().NoError(err)

			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			s.Require().NoError(err)

			assert.Equal(s.T(), c.expCode, resp.StatusCode)
			assert.Equal(s.T(), string(respBody), c.expRes)
		})
	}
}

func (s *ServerCurrencySuite) TestDeleteCurrency() {
	ctx := context.Background()

	testCases := []struct {
		name     string
		id       string
		mockFunc func()
		expRes   string
		expCode  int
	}{
		{
			name: "success",
			id:   uuid.New().String(),
			mockFunc: func() {
				s.mockCurrencySvc.On("DeleteCurrency", ctx, mock.Anything).
					Return(nil).Once()
			},
			expRes:  "",
			expCode: 204,
		},
		{
			name:     "invalid_id",
			id:       "123",
			mockFunc: func() {},
			expRes:   `{"code":400,"text":"invalid id format"}`,
			expCode:  400,
		},
		{
			name: "currency_not_found",
			id:   uuid.New().String(),
			mockFunc: func() {
				s.mockCurrencySvc.On("DeleteCurrency", ctx, mock.Anything).
					Return(entity.ErrEntityNotFound).Once()
			},
			expRes:  `{"code":404,"text":"Not Found"}`,
			expCode: 404,
		},
		{
			name: "svc_err",
			id:   uuid.New().String(),
			mockFunc: func() {
				s.mockCurrencySvc.On("DeleteCurrency", ctx, mock.Anything).
					Return(errors.New("")).Once()
			},
			expRes:  `{"code":500,"text":"Internal Server error"}`,
			expCode: 500,
		},
	}

	for _, c := range testCases {
		s.T().Run(c.name, func(t *testing.T) {
			app := fiber.New()
			app.Delete("/:id", s.srv.DeleteCurrency)

			c.mockFunc()

			req := httptest.NewRequest("DELETE", "/"+c.id, nil)

			resp, err := app.Test(req, 1)
			s.Require().NoError(err)

			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			s.Require().NoError(err)

			assert.Equal(s.T(), c.expCode, resp.StatusCode)
			assert.Equal(s.T(), string(respBody), c.expRes)
		})
	}
}

func (s *ServerCurrencySuite) TestConvert() {
	ctx := context.Background()

	testCases := []struct {
		name     string
		params   string
		mockFunc func()
		expRes   string
		expCode  int
	}{
		{
			name:   "success",
			params: "?from=USD&to=BTC&amount=1",
			mockFunc: func() {
				s.mockCurrencySvc.On("Convert", ctx, mock.Anything).
					Return(123.54, nil).Once()
			},
			expRes:  `{"course":123.54}`,
			expCode: 200,
		},
		{
			name:     "validation_err",
			params:   "?from=USD&to=BTC",
			mockFunc: func() {},
			expRes:   `{"code":400,"text":"Key: 'ConvertCurrencyReq.Amount' Error:Field validation for 'Amount' failed on the 'required' tag"}`,
			expCode:  400,
		},
		{
			name:   "svc_err",
			params: "?from=USD&to=BTC&amount=1",
			mockFunc: func() {
				s.mockCurrencySvc.On("Convert", ctx, mock.Anything).
					Return(float64(0), errors.New("")).Once()
			},
			expRes:  `{"code":400,"text":"Bad Request","businessCode":1}`,
			expCode: 400,
		},
	}

	for _, c := range testCases {
		s.T().Run(c.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/", s.srv.Convert)

			c.mockFunc()

			req := httptest.NewRequest("GET", "/"+c.params, nil)

			resp, err := app.Test(req, 10)
			s.Require().NoError(err)

			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			s.Require().NoError(err)

			assert.Equal(s.T(), c.expCode, resp.StatusCode)
			assert.Equal(s.T(), string(respBody), c.expRes)
		})
	}
}
