package npool

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/middlewares"
	"github.com/traefik/traefik/v2/pkg/tracing"
)

const (
	basicTypeName = "HeadersToBody"
)

type headersToBody struct {
	next        http.Handler
	name        string
	headerNames []string
}

// NewBasic creates a headersToBody middleware.
func NewHeadersToBody(ctx context.Context, next http.Handler, config dynamic.HeadersToBody, name string) (http.Handler, error) {
	middlewares.GetLogger(ctx, name, basicTypeName).Debug().Msg("Creating middleware")

	ctb := &headersToBody{
		name:        name,
		next:        next,
		headerNames: config.HeaderNames,
	}

	return ctb, nil
}

func (ctb *headersToBody) GetTracingInformation() (string, ext.SpanKindEnum) {
	return ctb.name, tracing.SpanKindNoneEnum
}

func (ctb *headersToBody) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger := middlewares.GetLogger(req.Context(), ctb.name, basicTypeName)

	myBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Warn().Msgf("Read body failed: %v", err)
		tracing.SetErrorWithEvent(req, "Read body failed")
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	bodyMap := map[string]interface{}{}
	if len(myBody) > 0 {
		err = json.Unmarshal(myBody, &bodyMap)
		if err != nil {
			logger.Warn().Msgf("Unmarshal body failed: %v", err)
			tracing.SetErrorWithEvent(req, "Unmarshal body failed")
			rw.WriteHeader(http.StatusForbidden)
			return
		}
	}

	for _, name := range ctb.headerNames {
		header := req.Header.Get(name)
		if header == "" {
			continue
		}

		bodyName := ""
		switch name {
		case authHeaderApp:
			bodyName = "AppID"
		case authHeaderUser:
			bodyName = "UserID"
		case authHeaderRole:
			bodyName = "Token"
		}

		if bodyName == "" {
			logger.Warn().Msgf("unexpected header to body")
			continue
		}

		bodyMap[bodyName] = header
	}

	myBody, err = json.Marshal(&bodyMap)
	if err != nil {
		logger.Warn().Msgf("Marshal body failed: %v", err)
		tracing.SetErrorWithEvent(req, "Marshal body failed")
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	req.Body = ioutil.NopCloser(strings.NewReader(string(myBody)))
	req.ContentLength = int64(len(myBody))

	logger.Debug().Msg("header parsed successed")
	ctb.next.ServeHTTP(rw, req)
}
