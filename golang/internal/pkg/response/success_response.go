package response

import "github.com/gofiber/fiber/v2"

type success struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(success{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}
