// Package router chi router provider
package router

import "github.com/go-chi/chi/v5"

// New chi router constructor
func New() *chi.Mux {

	r := chi.NewRouter()

	return r
}
