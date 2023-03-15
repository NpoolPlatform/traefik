package npool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/middlewares"
	"github.com/traefik/traefik/v2/pkg/tracing"

	"github.com/go-resty/resty/v2"
)

const (
	authTypeName = "RBACAuth"
	authHost     = "authing-gateway.kube-system.svc.cluster.local:50250"
)

type rbacAuth struct {
	next        http.Handler
	name        string
	headerNames []string
}

// NewRBAC creates a forward auth middleware.
func NewRBAC(ctx context.Context, next http.Handler, config dynamic.RBACAuth, name string) (http.Handler, error) {
	middlewares.GetLogger(ctx, name, authTypeName).Debug().Msg("Creating middleware")

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
	logger := middlewares.GetLogger(req.Context(), ra.name, authTypeName)

	userID := ""
	appID := ""
	userToken := ""

	ok := true
	for _, name := range ra.headerNames {
		header := req.Header.Get(name)
		if header == "" {
			logger.Warn().Msgf("fail get header %v", name)
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
		logger.Warn().Msgf("invalid app id")
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
		_, err = resty.New().R().
			SetBody(aReq).
			SetResult(&aResp).
			Post(fmt.Sprintf("http://%v/v1/auth/by/app/role/user", authHost))
	} else {
		_, err = resty.New().R().
			SetBody(aReq).
			SetResult(&aResp).
			Post(fmt.Sprintf("http://%v/v1/auth/by/app", authHost))
	}

	if err != nil {
		logger.Error().Msgf("fail auth by app: %v", err)
		ok = false
		goto lFail
	}

	if !aResp.Allowed {
		logger.Warn().Msgf("forbidden access")
		ok = false
	}

lFail:
	if !ok {
		logger.Warn().Msgf("header parse failed")
		tracing.SetErrorWithEvent(req, "header parse failed")
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	ra.next.ServeHTTP(rw, req)
}
