package npool

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/log"
	"github.com/traefik/traefik/v2/pkg/middlewares"
	"github.com/traefik/traefik/v2/pkg/tracing"

	"github.com/go-resty/resty/v2"

	"github.com/google/uuid"
)

const (
	opLogTypeName = "OpLog"
	opLogHost     = "oplog-gateway.kube-system.svc.cluster.local:50790"
)

type opLog struct {
	next http.Handler
	name string
}

// NewOpLog creates a forward oplog middleware
func NewOpLog(ctx context.Context, next http.Handler, config dynamic.OpLog, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, opLogTypeName)).Debug("Creating middleware")

	ol := &opLog{
		name: name,
		next: next,
	}

	return ol, nil
}

func (ol *opLog) GetTracingInformation() (string, ext.SpanKindEnum) {
	return ol.name, tracing.SpanKindNoneEnum
}

func (ol *opLog) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" && req.Method != "GET" {
		ol.next.ServeHTTP(rw, req)
		return
	}

	logger := log.FromContext(middlewares.GetLoggerCtx(req.Context(), ol.name, opLogTypeName))

	type opLogReq struct {
		AppID     string
		UserID    *string
		Path      string
		Method    string
		Arguments string
	}

	olq := &opLogReq{
		AppID:  uuid.UUID{}.String(),
		Path:   req.URL.String(),
		Method: "POST",
	}

	header := req.Header.Get("X-App-ID")
	if _, err := uuid.Parse(header); err != nil {
		logger.Warnf("Parse X-App-ID %v failed: %v", header, err)
		tracing.SetErrorWithEvent(req, "Parse X-App-ID failed")
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	olq.AppID = header

	header = req.Header.Get("X-User-ID")
	if _, err := uuid.Parse(header); err == nil {
		olq.UserID = &header
	}

	_body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Warnf("Read body failed: %v", err)
		tracing.SetErrorWithEvent(req, "Read body failed")
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	olq.Arguments = string(_body)

	type opLogResp struct {
		Info map[string]map[string]string
	}
	var resp *resty.Response

	resp, err = resty.
		New().
		R().
		SetBody(olq).
		SetResult(&opLogResp{}).
		Post(fmt.Sprintf("http://%v/v1/create/oplog", opLogHost))
	if err == nil {
		logger.Warnf("fail create oplog: %v", err)
	}

	req.Body = ioutil.NopCloser(strings.NewReader(string(_body)))
	buffer := bytes.NewBuffer(nil)
	_rw := newMultiWriter(buffer, rw)

	logger.Infof("oplog successed, url=%v, host=%v", req.URL, req.Host)
	ol.next.ServeHTTP(_rw, req)
	logger.Infof("oplog next done, url=%v, host=%v", req.URL, req.Host)

	olr := resp.Result().(*opLogResp)

	_olq, _ := json.Marshal(olq)
	_olr, _ := json.Marshal(olr)

	logger.Infof(
		"oplog serve done, url=%v, host=%v, req=%v, resp=%v, resp_len=%v, olq=%v, olr=%v",
		req.URL,
		req.Host,
		string(_body),
		buffer.String(),
		buffer.Len(),
		_olq,
		_olr,
	)
}

type multiWriter struct {
	bufWriter *bufio.Writer
	rw        http.ResponseWriter
	multi     io.Writer
}

func newMultiWriter(buffer *bytes.Buffer, rw http.ResponseWriter) http.ResponseWriter {
	bufWriter := bufio.NewWriter(buffer)
	multi := io.MultiWriter(bufWriter, rw)
	return &multiWriter{
		bufWriter: bufWriter,
		rw:        rw,
		multi:     multi,
	}
}

func (mw *multiWriter) Header() http.Header {
	return mw.rw.Header()
}

func (mw *multiWriter) Write(b []byte) (int, error) {
	return mw.multi.Write(b)
}

func (mw *multiWriter) WriteHeader(statusCode int) {
	mw.rw.WriteHeader(statusCode)
}
