package v1

import (
	"context"
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/veleton777/test_work_blum/internal/currency/v1/currency/entity"
	"github.com/veleton777/test_work_blum/internal/dto"
	"github.com/veleton777/test_work_blum/internal/pkg/httputil"
)

const codeCurrencyNotAllowedForConvert = 1

type CurrencyServer struct {
	currencySvc CurrencySvc
	validator   *validator.Validate
}

//go:generate mockery --name CurrencySvc
type CurrencySvc interface {
	CreateCurrency(ctx context.Context, currency dto.Currency) error
	UpdateCurrency(ctx context.Context, currency dto.Currency) error
	DeleteCurrency(ctx context.Context, id uuid.UUID) error

	Convert(ctx context.Context, course dto.ConvertCurrencyReq) (float64, error)
}

func NewCurrencyServer(currencySvc CurrencySvc) *CurrencyServer {
	return &CurrencyServer{
		currencySvc: currencySvc,
		validator:   validator.New(),
	}
}

// CreateCurrency godoc
//
//	@Summary		Create new currency
//	@Description	Create new currency
//	@Tags			currency
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.Currency	true	"CreateCurrencyDTO"
//	@Success		201
//	@Failure		400		{object}  httputil.HTTPError
//	@Router			/v1/currencies [post]
func (s *CurrencyServer) CreateCurrency(c *fiber.Ctx) error {
	var currency dto.Currency
	if err := json.Unmarshal(c.Body(), &currency); err != nil {
		return httputil.NewBadRequestErr(c, "invalid json body format") //nolint:wrapcheck
	}

	if err := s.validator.Struct(currency); err != nil {
		return httputil.NewBadRequestErr(c, err.Error()) //nolint:wrapcheck
	}

	if err := s.currencySvc.CreateCurrency(c.UserContext(), currency); err != nil {
		if errors.Is(err, entity.ErrCurrencyAlreadyExists) {
			return httputil.NewBadRequestErr(c, "currency already exists") //nolint:wrapcheck
		}

		return httputil.NewInternalServerErr(c) //nolint:wrapcheck
	}

	return httputil.NewCreatedResponse(c) //nolint:wrapcheck
}

// UpdateCurrency godoc
//
//		@Summary		Update currency
//		@Description	Update currency
//		@Tags			currency
//		@Accept			json
//		@Produce		json
//	    @Param          id   path string  true  "CurrencyID" Format(uuid)
//		@Param			payload	body		dto.Currency	true	"UpdateCurrencyDTO"
//		@Success		204
//		@Failure		400		{object}  httputil.HTTPError
//		@Failure		404		{object}  httputil.HTTPError
//		@Router			/v1/currencies/{id} [put]
func (s *CurrencyServer) UpdateCurrency(c *fiber.Ctx) error {
	id := c.Params("id")

	currencyID, err := uuid.Parse(id)
	if err != nil {
		return httputil.NewBadRequestErr(c, "invalid id format") //nolint:wrapcheck
	}

	var currency dto.Currency
	if err = json.Unmarshal(c.Body(), &currency); err != nil {
		return httputil.NewBadRequestErr(c, "invalid json body format") //nolint:wrapcheck
	}

	if err = s.validator.Struct(currency); err != nil {
		return httputil.NewBadRequestErr(c, err.Error()) //nolint:wrapcheck
	}

	currency.ID = currencyID

	if err = s.currencySvc.UpdateCurrency(c.UserContext(), currency); err != nil {
		if errors.Is(err, entity.ErrEntityNotFound) {
			return httputil.NewNotFoundErr(c) //nolint:wrapcheck
		}

		return httputil.NewInternalServerErr(c) //nolint:wrapcheck
	}

	return httputil.NewNoContentResponse(c) //nolint:wrapcheck
}

// DeleteCurrency godoc
//
//	@Summary		Delete currency
//	@Description	Delete currency
//	@Tags			currency
//	@Accept			json
//	@Produce		json
//	@Param          id   path string  true  "CurrencyID" Format(uuid)
//	@Success		204
//	@Failure		400		{object}  httputil.HTTPError
//	@Failure		404		{object}  httputil.HTTPError
//	@Router			/v1/currencies/{id} [delete]
func (s *CurrencyServer) DeleteCurrency(c *fiber.Ctx) error {
	id := c.Params("id")

	currencyID, err := uuid.Parse(id)
	if err != nil {
		return httputil.NewBadRequestErr(c, "invalid id format") //nolint:wrapcheck
	}

	if err = s.currencySvc.DeleteCurrency(c.UserContext(), currencyID); err != nil {
		if errors.Is(err, entity.ErrEntityNotFound) {
			return httputil.NewNotFoundErr(c) //nolint:wrapcheck
		}

		return httputil.NewInternalServerErr(c) //nolint:wrapcheck
	}

	return httputil.NewNoContentResponse(c) //nolint:wrapcheck
}

// Convert godoc
//
//		@Summary		Convert course for currencies
//		@Description	Convert course for currencies
//		@Tags			currency
//		@Accept			json
//		@Produce		json
//		@Param			payload	query		dto.ConvertCurrencyReq	true	"ConvertCurrencyReq"
//		@Success		200		{object}  dto.ConvertCurrencyResp
//		@Failure		400		{object}  httputil.HTTPError
//	    @Router			/v1/currencies/convert [get]
func (s *CurrencyServer) Convert(c *fiber.Ctx) error {
	var req dto.ConvertCurrencyReq
	if err := c.QueryParser(&req); err != nil {
		return httputil.NewBadRequestErr(c, "invalid query params") //nolint:wrapcheck
	}

	if err := s.validator.Struct(req); err != nil {
		return httputil.NewBadRequestErr(c, err.Error()) //nolint:wrapcheck
	}

	res, err := s.currencySvc.Convert(c.UserContext(), req)
	if err != nil {
		return httputil.NewBusinessErr(c, codeCurrencyNotAllowedForConvert) //nolint:wrapcheck
	}

	resp := dto.ConvertCurrencyResp{Course: res}

	return c.JSON(resp) //nolint:wrapcheck
}
