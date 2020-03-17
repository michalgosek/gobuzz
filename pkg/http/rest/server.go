package rest

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuzz/pkg/domain/adding"
	"github.com/gobuzz/pkg/domain/responding"
)

type server struct {
	router *chi.Mux
}

// ServHandler creates server handler and returns registered router
func ServHandler(a adding.Service, r responding.Service) *chi.Mux {
	s := newServer(a, r)
	return s.router
}

func newServer(a adding.Service, r responding.Service) *server {
	s := &server{
		router: chi.NewRouter(),
	}
	s.router.Use(middleware.Logger)
	s.routes(a, r)
	return s
}
