package middlewares

import (
	"net/http"
)

type Middlewares struct {
	Rbac func(next http.HandlerFunc, permission string) http.HandlerFunc
	Auth func(next http.HandlerFunc) http.HandlerFunc
}

type MiddelwaresBuilder struct {
	middlewares *Middlewares
}

func NewMiddlewaresBuilder() *MiddelwaresBuilder {
	return &MiddelwaresBuilder{
		middlewares: &Middlewares{},
	}
}

func (b *MiddelwaresBuilder) WithAuth() *MiddelwaresBuilder {
	b.middlewares.Auth = b.middlewares.authFunc
	return b
}

func (b *MiddelwaresBuilder) WithRBAC() *MiddelwaresBuilder {
	b.middlewares.Rbac = b.middlewares.rbacFunc
	return b
}

func (b *MiddelwaresBuilder) Build() *Middlewares {
	return b.middlewares
}
