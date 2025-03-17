package npool

import (
	"context"
	"fmt"
	"io/ioutil"
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
	billingTypeName = "Billing"
	billingHost     = "billing-gateway.kube-system.svc.cluster.local:50900"
)

type billing struct {
	next        http.Handler
	name        string
	headerNames []string
}

// NewRBAC creates a forward auth middleware.
func NewBilling(ctx context.Context, next http.Handler, config dynamic.Billing, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, billingTypeName)).Debug("Creating middleware")

	ra := &billing{
		name: name,
		next: next,
	}

	return ra, nil
}

func (ra *billing) GetTracingInformation() (string, ext.SpanKindEnum) {
	return ra.name, tracing.SpanKindNoneEnum
}

func (ra *billing) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger := log.FromContext(middlewares.GetLoggerCtx(req.Context(), ra.name, authTypeName))
	var userID *string
	var appID string

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
		}
	}

	if appID == "" {
		logger.Warnf("invalid app id")
		ok = false
	}

	var err error

	type billingReq struct {
		AppID  string
		UserID *string
		Path   string
		ReqMsg string
	}

	type billingResp struct {
		Allow bool // Here info is allowed
	}

	_body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Warnf("Read body failed: %v", err)
		tracing.SetErrorWithEvent(req, "Read body failed")
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	aReq := billingReq{
		AppID:  appID,
		UserID: userID,
		Path:   req.URL.String(),
		ReqMsg: string(_body),
	}
	var aResp *billingResp
	var resp *resty.Response

	if !ok {
		goto lFail
	}

	resp, err = resty.
		New().
		R().
		SetBody(aReq).
		SetResult(&billingResp{}).
		Post(fmt.Sprintf("http://%v/v1/user/calculate/charge", billingHost))
	if err != nil {
		logger.Errorf("fail auth: %v", err)
		ok = false
		goto lFail
	}

	aResp = resp.Result().(*billingResp)
	if !aResp.Allow {
		logger.Warnf("forbidden access: %v", resp)
		ok = false
	}

lFail:
	if !ok {
		logger.Warnf("billing failed")
		tracing.SetErrorWithEvent(req, "billing failed")
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	ra.next.ServeHTTP(rw, req)
}
