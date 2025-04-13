package exception

import (
	"github.com/gofiber/fiber/v2"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
)

func defaultStackTraceHandler(_ *fiber.Ctx, e interface{}) {
}

func configDefault(config ...fiberRecover.Config) fiberRecover.Config {
	if len(config) < 1 {
		return fiberRecover.ConfigDefault
	}

	cfg := config[0]
	if cfg.EnableStackTrace && cfg.StackTraceHandler == nil {
		cfg.StackTraceHandler = defaultStackTraceHandler
	}

	return cfg
}

func Recover(config ...fiberRecover.Config) fiber.Handler {
	cfg := configDefault(config...)

	return func(c *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		defer func() {
			if r := recover(); r != nil {
				if cfg.StackTraceHandler != nil {
					cfg.StackTraceHandler(c, r)
				}

				if _, ok := r.(error); ok {
				}

				c.SendStatus(fiber.StatusInternalServerError)
			}
		}()

		return c.Next()
	}
}
