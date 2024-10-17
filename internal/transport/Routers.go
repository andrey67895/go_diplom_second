package transport

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func GetRoutersGophermart() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RealIP, middleware.Recoverer, middleware.Logger)
	return r
}
