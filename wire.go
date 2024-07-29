//go:build wireinject
// +build wireinject

package main

import (
	"github.com/devmizumizurice/go-jwt-graphql/initializers"
	"github.com/devmizumizurice/go-jwt-graphql/repositories"
	"github.com/devmizumizurice/go-jwt-graphql/services"
	"github.com/google/wire"
)

func InitializeService() services.ServicesInterface {
	wire.Build(
		services.NewAuthService,
		services.NewUserService,
		repositories.NewUserRepository,
		initializers.GetDB,
		services.NewServices,
	)
	return nil
}
