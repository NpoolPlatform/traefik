package auth

import (
	"context"
	"net/http"

	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/log"
	"github.com/traefik/traefik/v2/pkg/middlewares"
)

type rbacAuth struct {
}

// NewRBAC creates a forward auth middleware.
func NewRBAC(ctx context.Context, next http.Handler, config dynamic.RBACAuth, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, forwardedTypeName)).Debug("Creating middleware")

	ra := &rbacAuth{}

	return ra, nil
}

func (ra *rbacAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// TODO: check app exist
	// TODO: check user exist
	// TODO: check user login
	// TODO: check user session
	// TODO: check user permission
	fa.next.ServeHTTP(rw, req)
}
