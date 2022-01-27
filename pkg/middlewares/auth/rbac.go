package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/log"
	"github.com/traefik/traefik/v2/pkg/middlewares"
	"github.com/traefik/traefik/v2/pkg/tracing"

	"github.com/go-resty/resty/v2"
)

const (
	authTypeName   = "RBACAuth"
	authHeaderApp  = "X-App-ID"
	authHeaderUser = "X-User-ID"
	authHeaderRole = "X-App-Login-Token"
	authHost       = "authing-gateway.kube-system.svc.cluster.local:50250"
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

	var err error

	type authReq struct {
		AppID    string
		UserID   string
		Token    string
		Resource string
		Method   string
	}

	type authResp struct {
		Allowed bool
	}
	aReq := authReq{
		AppID:    appID,
		UserID:   userID,
		Token:    userToken,
		Resource: req.URL.String(),
		Method:   req.Method,
	}
	aResp := authResp{}

	if !ok {
		goto lFail
	}

	if userID != "" {
		if userToken == "" {
			ok = false
			goto lFail
		}
		// TODO: user authorize
		_, err = resty.New().R().
			SetBody(aReq).
			SetResult(&aResp).
			Post(fmt.Sprintf("http://%v/v1/auth/by/app/role/user"))
	} else {
		_, err = resty.New().R().
			SetBody(aReq).
			SetResult(&aResp).
			Post(fmt.Sprintf("http://%v/v1/auth/by/app"))
	}

	if err != nil {
		logger.Errorf("fail auth by app: %v", err)
		ok = false
		goto lFail
	}

	if !aResp.Allowed {
		logger.Warnf("forbidden access")
		ok = false
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
