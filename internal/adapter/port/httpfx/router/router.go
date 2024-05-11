// Package router chi router provider
package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"moul.io/chizap"

	"go.uber.org/zap"

	dpkg "github.com/structx/go-dpkg/adapter/port/http/controller"

	"githhub.com/structx/ddns/internal/adapter/port/httpfx/controller"
	"githhub.com/structx/ddns/internal/core/domain"
)

// New chi router constructor
func New(logger *zap.Logger, ddns domain.DDNS) *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(chizap.New(logger, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	cc := []interface{}{
		dpkg.NewBundle(logger),
		controller.NewRecords(logger, ddns),
	}

	v1 := chi.NewRouter()

	for _, c1 := range cc {

		if c, ok := c1.(dpkg.V0); ok {
			h := c.RegisterRoutesV0()
			r.Mount("/", h)
		}

		if c, ok := c1.(dpkg.V1); ok {
			h := c.RegisterRoutesV1()
			v1.Mount("/", h)
		}
	}

	r.Mount("/api/v1", v1)

	return r
}
