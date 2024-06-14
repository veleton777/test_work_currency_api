//nolint:godot
package server

import (
	"github.com/gofiber/fiber/v2"
)

// @title Swagger Currency API
// @version 1.0
// @description This is a server Currency API

// @host localhost:8080
// @BasePath /api
func (s *API) routes(app *fiber.App) {
	api := app.Group("api")
	api.Post("/v1/currencies", s.currencyServer.CreateCurrency)
	api.Put("/v1/currencies/:id", s.currencyServer.UpdateCurrency)
	api.Delete("/v1/currencies/:id", s.currencyServer.DeleteCurrency)

	api.Get("/v1/currencies/convert", s.currencyServer.Convert)
}
