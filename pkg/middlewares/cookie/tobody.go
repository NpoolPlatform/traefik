package cookie

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/log"
	"github.com/traefik/traefik/v2/pkg/middlewares"
	"github.com/traefik/traefik/v2/pkg/tracing"
)

const (
	basicTypeName = "CookiesToBody"
)

type cookieToBody struct {
	next	    http.Handler
	name        string
	cookieNames []string
}

// NewBasic creates a cookieToBody middleware.
func New(ctx context.Context, next http.Handler, config dynamic.CookiesToBody, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, basicTypeName)).Debug("Creating middleware")

	ctb := &cookieToBody{
		name:	     name,
		next:	     next,
		cookieNames: config.CookieNames,
	}

	return ctb, nil
}

func (ctb *cookieToBody) GetTracingInformation() (string, ext.SpanKindEnum) {
	return ctb.name, tracing.SpanKindNoneEnum
}

func (ctb *cookieToBody) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger := log.FromContext(middlewares.GetLoggerCtx(req.Context(), ctb.name, basicTypeName))

	myBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Warnf("Read body failed: %v", err)
		tracing.SetErrorWithEvent(req, "Read body failed")
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	bodyMap := map[string]interface{}{}
	if len(myBody) > 0 {
		err = json.Unmarshal(myBody, &bodyMap)
		if err != nil {
			logger.Warnf("Unmarshal body failed: %v", err)
			tracing.SetErrorWithEvent(req, "Unmarshal body failed")
			rw.WriteHeader(http.StatusForbidden)
			return
		}
	}

	ok := true
	for _, name := range ctb.cookieNames {
		cookie, err := req.Cookie(name)
		if err != nil {
			logger.Warnf("Cookie %v error %v", name, err)
			ok = false
			continue
		}

		bodyMap[cookie.Name] = cookie.Value
	}

	if !ok {
		logger.Warnf("Cookie parse failed")
		tracing.SetErrorWithEvent(req, "Cookie parse failed")
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	myBody, err = json.Marshal(&bodyMap)
	if err != nil {
		logger.Warnf("Marshal body failed: %v", err)
		tracing.SetErrorWithEvent(req, "Marshal body failed")
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	req.Body = ioutil.NopCloser(strings.NewReader(string(myBody)))
	req.ContentLength = int64(len(myBody))

	logger.Debug("Cookie parsed successed")
	ctb.next.ServeHTTP(rw, req)
}
