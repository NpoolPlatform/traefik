package cookie

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/log"
	"github.com/traefik/traefik/v2/pkg/middlewares"
	"github.com/traefik/traefik/v2/pkg/tracing"

	"github.com/google/uuid"
)

const (
	basicTypeName = "CookieToBody"
)

type cookieToBody struct {
	cookieNames []string
}

// NewBasic creates a cookieToBody middleware.
func NewCookieToBody(ctx context.Context, next http.Handler, config dynamic.CookieToBody, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, basicTypeName)).Debug("Creating middleware")

	ctb := &cookieToBody{
		cookieNames: config.CookieNames,
	}

	return ctb, nil
}

func (ctb *cookieToBody) GetTracingInformation() (string, ext.SpanKindEnum) {
	return b.name, tracing.SpanKindNoneEnum
}

func (ctb *cookieToBody) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger := log.FromContext(middlewares.GetLoggerCtx(req.Context(), b.name, basicTypeName))

	myBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Debug("Read body failed: %v", err)
		tracing.SetErrorWithEvent(req, "Read body failed")
		return
	}

	bodyMap := map[string]interface{}{}
	err = json.Unmarshal(myBody, &bodyMap)
	if err != nil {
		logger.Debug("Unmarshal body failed: %v", err)
		tracing.SetErrorWithEvent(req, "Unmarshal body failed")
		return
	}

	ok := true
	for name := range ctb.cookieNames {
		cookie, err := r.Cookie(name)
		if err != nil {
			logger.Debug("Cookie %v error %v", name, err)
			ok = false
		}

		bodyMap[cookie.Name] = cookie.Value
	}

	if !ok {
		logger.Debug("Cookie parse failed")
		tracing.SetErrorWithEvent(req, "Cookie parse failed")
		return
	}

	myBody, err = json.Marshal(&bodyMap)
	if err != nil {
		logger.Debug("Marshal body failed: %v", err)
		tracing.SetErrorWithEvent(req, "Marshal body failed")
		return
	}

	r.Body = ioutil.NopCloser(strings.NewReader(string(myBody)))
	r.ContentLength = int64(len(b))

	logger.Debug("Cookie parsed successed")
	b.next.ServeHTTP(rw, req)
}
