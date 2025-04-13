package exception

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kyh0703/template/internal/pkg/logger"
	"github.com/kyh0703/template/internal/pkg/response"
)

func errorResponse(err error) (int, response.Error) {
	logger.Debug("error-response", err)

	if ce, ok := err.(*Error); ok {
		return fiber.StatusBadRequest, response.Error{
			Code:    ce.Code,
			Message: ce.Message,
			Detail:  ce.Detail,
		}
	}

	if fe, ok := err.(*fiber.Error); ok {
		return fe.Code, response.Error{
			Code:    CodeServerInternal,
			Message: fe.Message,
			Detail:  "",
		}
	}

	return 500, response.Error{
		Message: "Internal Server Error",
		Detail:  "",
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	errCode, errBody := errorResponse(err)
	return c.Status(errCode).JSON(errBody)
}
