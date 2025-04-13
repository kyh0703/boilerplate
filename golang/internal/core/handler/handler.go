package handler

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Mapper struct {
	Method  string
	Path    string
	Handler []fiber.Handler
}

func Mapping(method, path string, handler ...fiber.Handler) Mapper {
	return Mapper{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}

type Handler interface {
	Table() []Mapper
}

func AsHandler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Handler)),
		fx.ResultTags(`group:"handlers"`),
	)
}
