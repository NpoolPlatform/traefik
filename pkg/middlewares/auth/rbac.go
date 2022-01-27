package auth

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/log"
	"github.com/traefik/traefik/v2/pkg/middlewares"
	"github.com/traefik/traefik/v2/pkg/tracing"
)

const (
	authTypeName   = "RBACAuth"
	authHeaderApp  = "X-App-ID"
	authHeaderUser = "X-User-ID"
	authHeaderRole = "X-App-Login-Token"
)

type rbacAuth struct {
	next        http.Handler
	name        string
	headerNames []string
}

// NewRBAC creates a forward auth middleware.
func NewRBAC(ctx context.Context, next http.Handler, config dynamic.RBACAuth, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, forwardedTypeName)).Debug("Creating middleware")

	ra := &rbacAuth{
		name:        name,
		next:        next,
		headerNames: config.HeaderNames,
	}

	return ra, nil
}

func (ra *rbacAuth) GetTracingInformation() (string, ext.SpanKindEnum) {
	return ra.name, tracing.SpanKindNoneEnum
}

func (ra *rbacAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// TODO: check app exist
	// TODO: check user exist
	// TODO: check user login
	// TODO: check user session
	// TODO: check user permission
	logger := log.FromContext(middlewares.GetLoggerCtx(req.Context(), ra.name, authTypeName))

	userID := ""
	appID := ""
	userToken := ""

	ok := true
	for _, name := range ra.headerNames {
		header := req.Header.Get(name)
		if header == "" {
			logger.Warnf("fail get header %v", name)
			ok = false
			continue
		}

		switch name {
		case authHeaderApp:
			appID = req.Header.Get(authHeaderApp)
		case authHeaderUser:
			userID = req.Header.Get(authHeaderUser)
		case authHeaderRole:
			userToken = req.Header.Get(authHeaderRole)
		}
	}

	if appID == "" {
		logger.Warnf("invalid app id")
		ok = false
	}

	if !ok {
		goto lFail
	}

	if userID != "" {
		if userToken == "" {
			ok = false
			goto lFail
		}
		// TODO: user authorize
	} else {
		// TODO: app authorize
	}

lFail:
	if !ok {
		logger.Warnf("header parse failed")
		tracing.SetErrorWithEvent(req, "header parse failed")
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	ra.next.ServeHTTP(rw, req)
}
