package web

import (
	"net/http"

	"github.com/jmjp/go-rbac/internal/adapters/deliveries/web/handlers"
	repositories "github.com/jmjp/go-rbac/internal/adapters/repositories/mongo"
	"github.com/jmjp/go-rbac/internal/core/usecases"
)

func (s *Server) healthRoutes() {
	s.router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	s.router.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	})
}

func (s *Server) authRoutes() {
	userRepo := repositories.NewUserMongoRepository(s.database)
	otpRepo := repositories.NewOTPMongoRepository(s.database)
	sessionRepo := repositories.NewSessionMongoRepository(s.database)

	usecases := usecases.NewAuthUseCase(userRepo, otpRepo, sessionRepo)

	handler := handlers.NewAuthHandler(usecases)

	s.router.HandleFunc("POST /auth/login", handler.Login)
	s.router.HandleFunc("GET /auth/verify", handler.Verify)
	s.router.HandleFunc("POST /auth/refresh", handler.Refresh)
	s.router.HandleFunc("GET /auth/sessions", s.middlewares.Auth(handler.Sessions))
	s.router.HandleFunc("DELETE /auth/logout", s.middlewares.Auth(handler.Logout))

}

func (s *Server) teamsRoutes() {
	userRepo := repositories.NewUserMongoRepository(s.database)
	teamRepo := repositories.NewTeamMongoRepository(s.database)

	usecases := usecases.NewTeamUsecase(userRepo, teamRepo)

	handler := handlers.NewTeamHandler(usecases)

	s.router.HandleFunc("POST /teams", s.middlewares.Auth(handler.Create))
	s.router.HandleFunc("PUT /teams/{teamId}", s.middlewares.Auth(handler.Update))
	s.router.HandleFunc("DELETE /teams/{teamId}", s.middlewares.Auth(s.middlewares.Rbac(handler.Delete, "team::delete::*")))
}
