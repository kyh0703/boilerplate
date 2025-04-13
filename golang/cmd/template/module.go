package main

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"main",
	fx.Provide(
		fx.Annotate(NewFiber, fx.ParamTags(`group:"handlers"`)),
	),
)
