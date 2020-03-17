package rest

import (
	"github.com/go-chi/chi"
	"github.com/gobuzz/pkg/domain/adding"
	"github.com/gobuzz/pkg/domain/responding"
	"github.com/gobuzz/pkg/http/rest/handlers"
)

func (s *server) routes(adder adding.Service, respsr responding.Service) {

	s.router.Route("/api/fetcher", func(r chi.Router) {
		// r.Get("/", s.handleRequestsFetch())
		r.Post("/", handlers.HandleFetchCreate(adder, respsr))

		// r.Route("/{id}", func(r chi.Router) {
		// 	r.Get("/", s.handleRequestFetch())
		// 	r.Put("/", s.handleRequestUpdate())
		// })

	})

}
