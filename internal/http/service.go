package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var _ http.Handler = (*Service)(nil)

type Service struct {
	m *chi.Mux
}

// ServeHTTP implements http.Handler.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

func NewService() *Service {
	s := &Service{
		m: chi.NewMux(),
	}
	s.routes()
	return s
}

func (s *Service) routes() {
	s.m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}

		fmt.Fprintf(w, "Ciao, %s!", name)
	})
}
