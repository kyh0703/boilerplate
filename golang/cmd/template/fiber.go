package main

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/kyh0703/template/internal/core/handler"
	"github.com/kyh0703/template/internal/pkg/exception"
	"github.com/kyh0703/template/internal/pkg/logger"
)

func NewFiber(handlers ...handler.Handler) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "flow",
		ServerHeader: "flow",
		Prefork:      false,
		UnescapePath: true,
		ErrorHandler: exception.ErrorHandler,
	})

	app.Get("/swagger/*", swagger.HandlerDefault)
	app = setupHandlers(app, handlers...)
	app = setupMiddleware(app)
	return app
}

func setupMiddleware(app *fiber.App) *fiber.App {
	app.Use(cors.New())
	app.Use(exception.Recover())
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger.Logger.Desugar(),
	}))
	app.Use(pprof.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	return app
}

func setupHandlers(app *fiber.App, handlers ...handler.Handler) *fiber.App {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	for _, h := range handlers {
		for _, m := range h.Table() {
			v1.Add(m.Method, m.Path, m.Handler...)
		}
	}

	return app
}
