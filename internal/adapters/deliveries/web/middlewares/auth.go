package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/jmjp/go-rbac/internal/adapters/deliveries/tokens"
)

func (b *Middlewares) authFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if len(authorization) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tks := strings.Replace(authorization, "Bearer ", "", 1)
		tks = strings.Trim(tks, " ")
		if len(tks) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		payload, err := tokens.ValidatePasseto(tks)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		ctx := context.WithValue(r.Context(), "logged_user", payload)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
