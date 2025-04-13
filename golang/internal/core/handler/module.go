package handler

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"handler",
	fx.Provide(
		AsHandler(NewAuthHandler),
		AsHandler(NewEdgeHandler),
		AsHandler(NewNodeHandler),
		AsHandler(NewFlowHandler),
		AsHandler(NewProjectHandler),
		AsHandler(NewUserHandler),
	),
)
