package handler

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"handler",
	fx.Provide(
		AsHandler(NewAuthHandler),
		AsHandler(NewPostsHandler),
		AsHandler(NewUsersHandler),
	),
)
