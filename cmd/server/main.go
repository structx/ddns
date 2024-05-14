// Package main server entrypoint
package main

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/multierr"

	"github.com/structx/go-dpkg/adapter/logging"
	"github.com/structx/go-dpkg/adapter/port/http/serverfx"
	"github.com/structx/go-dpkg/adapter/setup"
	dpkg "github.com/structx/go-dpkg/domain"
	"github.com/structx/go-dpkg/util/decode"

	"github.com/structx/ddns/internal/adapter/port/httpfx/router"
	"github.com/structx/ddns/internal/adapter/port/rpcfx"
	"github.com/structx/ddns/internal/core/domain"
	"github.com/structx/ddns/internal/core/service"
)

func main() {
	fx.New(
		fx.Provide(context.TODO),
		fx.Provide(fx.Annotate(setup.New, fx.As(new(dpkg.Config)))),
		fx.Invoke(decode.ConfigFromEnv),
		fx.Provide(logging.New),
		fx.Provide(fx.Annotate(service.NewDDNS, fx.As(new(domain.DDNS)))),
		fx.Provide(fx.Annotate(router.New, fx.As(new(http.Handler)))),
		fx.Provide(rpcfx.NewGRPCServer),
		fx.Provide(serverfx.New),
		fx.Invoke(registerHooks),
	).Run()
}

func registerHooks(lc fx.Lifecycle, s1 *http.Server) error {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {

				var result error

				go func() {
					if err := s1.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						result = multierr.Append(result, fmt.Errorf("failed to start http server %v", err))
					}
				}()

				return result
			},
			OnStop: func(ctx context.Context) error {

				var result error

				err := s1.Shutdown(ctx)
				if err != nil {
					result = multierr.Append(result, fmt.Errorf("unable to shutdown http server %v", err))
				}

				return result
			},
		},
	)
	return nil
}
