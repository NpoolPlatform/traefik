package grpcweb

import (
	"context"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/middlewares"
)

const typeName = "grpc-web"

// New builds a new gRPC web request converter.
func New(ctx context.Context, next http.Handler, config dynamic.GrpcWeb, name string) http.Handler {
	middlewares.GetLogger(ctx, name, typeName).Debug().Msg("Creating middleware")

	return grpcweb.WrapHandler(next, grpcweb.WithCorsForRegisteredEndpointsOnly(false), grpcweb.WithOriginFunc(func(origin string) bool {
		middlewares.GetLogger(ctx, name, typeName).Debug().Msgf("GrpcWeb allowOrigins %v, origin %v", config.AllowOrigins, origin)
		for _, originCfg := range config.AllowOrigins {
			if originCfg == "*" || originCfg == origin {
				return true
			}
		}
		return false
	}))
}
