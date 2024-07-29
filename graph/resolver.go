package graph

import "github.com/devmizumizurice/go-jwt-graphql/services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Srv services.ServicesInterface
}
