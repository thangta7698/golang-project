package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"go-training-system/internal/service"
)

type Resolver struct {
	UserService service.UserService
	JWTSecret   string
}
