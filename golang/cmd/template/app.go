package main

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/kyh0703/template/configs"
)

type app struct {
	config *configs.Config
	fiber  *fiber.App
}

func NewApp(config *configs.Config, fiber *fiber.App) *app {
	return &app{
		config: config,
		fiber:  fiber,
	}
}

func (a *app) Run(ctx context.Context) error {
	a.fiber.Listen(":" + a.config.Server.Port)
	return nil
}

func (a *app) Stop(ctx context.Context) error {
	a.fiber.Shutdown()
	return nil
}
