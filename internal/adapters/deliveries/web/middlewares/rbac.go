package middlewares

import (
	"net/http"

	"github.com/jmjp/go-rbac/internal/adapters/deliveries/tokens"
	"github.com/jmjp/go-rbac/pkg/rbac"
)

func (b *Middlewares) rbacFunc(next http.HandlerFunc, permission string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rolesFile := rbac.NewRBACBuilder().WithRolesFromFile("./config/roles.json").Build()
		userId := r.Context().Value("logged_user").(*tokens.Payload)
		teamId := r.PathValue("teamId")
		for _, team := range userId.Teams {
			if team.TeamID == teamId && rolesFile.HasPermission(permission, team.Role) {
				next(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
	}
}
