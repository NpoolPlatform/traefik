package npool

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

	"github.com/google/uuid"
)

const (
	authTypeName = "RBACAuth"
	authHost     = "appuser-gateway.kube-system.svc.cluster.local:50500"
)

type rbacAuth struct {
	next        http.Handler
	name        string
	headerNames []string
}

// NewRBAC creates a forward auth middleware.
func NewRBAC(ctx context.Context, next http.Handler, config dynamic.RBACAuth, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, authTypeName)).Debug("Creating middleware")

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
	logger := log.FromContext(middlewares.GetLoggerCtx(req.Context(), ra.name, authTypeName))

	var userID *string
	var appID string
	var userToken *string

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
			_userID := req.Header.Get(authHeaderUser)
			if _, err := uuid.Parse(_userID); err == nil {
				userID = &_userID
			}
		case authHeaderRole:
			_userToken := req.Header.Get(authHeaderRole)
			userToken = &_userToken
		}
	}

	if appID == "" {
		logger.Warnf("invalid app id")
		ok = false
	}
	if userID != nil && userToken == nil {
		logger.Warnf("invalid userid & usertoken")
		ok = false
	}

	var err error

	type authReq struct {
		AppID    string
		UserID   *string
		Token    *string
		Resource string
		Method   string
	}

	type authResp struct {
		Info bool // Here info is allowed
	}
	aReq := authReq{
		AppID:    appID,
		UserID:   userID,
		Token:    userToken,
		Resource: req.URL.String(),
		Method:   req.Method,
	}
	var aResp *authResp
	var resp *resty.Response

	if !ok {
		goto lFail
	}

	resp, err = resty.
		New().
		R().
		SetBody(aReq).
		SetResult(&authResp{}).
		Post(fmt.Sprintf("http://%v/v1/authenticate", authHost))
	if err != nil {
		logger.Errorf("fail auth: %v", err)
		ok = false
		goto lFail
	}

	aResp = resp.Result().(*authResp)
	if !aResp.Info {
		logger.Warnf("forbidden access: %v", resp)
		ok = false
	}

lFail:
	if !ok {
		logger.Warnf("authorize failed")
		tracing.SetErrorWithEvent(req, "authorize failed")
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	ra.next.ServeHTTP(rw, req)
}
