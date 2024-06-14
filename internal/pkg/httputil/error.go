package httputil

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type HTTPError struct {
	Code         int    `json:"code"`
	Text         string `json:"text"`
	BusinessCode int    `json:"businessCode,omitempty"`
}

func NewNotFoundErr(ctx *fiber.Ctx) error {
	resp := HTTPError{
		Code:         fiber.StatusNotFound,
		Text:         "Not Found",
		BusinessCode: 0,
	}

	ctx.Status(fiber.StatusNotFound)

	if err := ctx.JSON(resp); err != nil {
		return errors.Wrap(err, "write json resp")
	}

	return nil
}

func NewBadRequestErr(ctx *fiber.Ctx, msg string) error {
	resp := HTTPError{
		Code:         fiber.StatusBadRequest,
		Text:         "Bad Request",
		BusinessCode: 0,
	}

	if msg != "" {
		resp.Text = msg
	}

	ctx.Status(fiber.StatusBadRequest)

	if err := ctx.JSON(resp); err != nil {
		return errors.Wrap(err, "write json resp")
	}

	return nil
}

func NewBusinessErr(ctx *fiber.Ctx, businessCode int) error {
	resp := HTTPError{
		Code:         fiber.StatusBadRequest,
		Text:         "Bad Request",
		BusinessCode: businessCode,
	}

	ctx.Status(fiber.StatusBadRequest)

	if err := ctx.JSON(resp); err != nil {
		return errors.Wrap(err, "write json resp")
	}

	return nil
}

func NewInternalServerErr(ctx *fiber.Ctx) error {
	resp := HTTPError{
		Code:         fiber.StatusInternalServerError,
		Text:         "Internal Server error",
		BusinessCode: 0,
	}

	ctx.Status(fiber.StatusInternalServerError)

	if err := ctx.JSON(resp); err != nil {
		return errors.Wrap(err, "write json resp")
	}

	return nil
}

func NewNoContentResponse(ctx *fiber.Ctx) error {
	ctx.Status(fiber.StatusNoContent)

	if err := ctx.SendString(""); err != nil {
		return errors.Wrap(err, "write empty resp")
	}

	return nil
}

func NewCreatedResponse(ctx *fiber.Ctx) error {
	ctx.Status(fiber.StatusCreated)

	if err := ctx.SendString(""); err != nil {
		return errors.Wrap(err, "write empty resp")
	}

	return nil
}
