package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmjp/go-rbac/internal/adapters/deliveries/web/middlewares"

	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	port        int
	router      *http.ServeMux
	database    *mongo.Database
	middlewares *middlewares.Middlewares
}

type ServerBuilder struct {
	server *Server
}

func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		server: &Server{
			router: http.NewServeMux(),
		},
	}
}

func (b *ServerBuilder) WithPort(port int) *ServerBuilder {
	b.server.port = port
	return b
}

func (b *ServerBuilder) WithDB(db *mongo.Client) *ServerBuilder {
	b.server.database = db.Database("zoops")
	return b
}

func (b *ServerBuilder) WithMiddlewares() *ServerBuilder {
	b.server.middlewares = middlewares.NewMiddlewaresBuilder().WithAuth().WithRBAC().Build()
	return b
}

// Build returns the built Server.
//
// No parameters.
// Returns a pointer to a Server.
func (b *ServerBuilder) Build() *Server {
	return b.server
}

// Start starts the HTTP server and listens for incoming requests.
//
// No parameters.
// No return value.
func (s *Server) Start() {
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.router,
	}

	s.healthRoutes()
	s.authRoutes()
	s.teamsRoutes()

	go func() {
		fmt.Printf("Starting HTTP server on port %d\n", s.port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("HTTP server error: %v", err)
			panic(err)
		}
		fmt.Printf("HTTP server stopping serving connections\n")
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-ch
	fmt.Printf("Graceful shutdown HTTP server\n")
	ctx, shuwdown := context.WithTimeout(context.Background(), time.Minute*10)
	defer shuwdown()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("HTTP server error: %v", err)
		panic(err)
	}
	fmt.Printf("HTTP server shutdown\n")
}
