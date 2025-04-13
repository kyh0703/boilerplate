package service

import (
	"github.com/kyh0703/template/internal/core/service/auth"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"service",
	fx.Provide(
		auth.NewAuthService,
	),
)
