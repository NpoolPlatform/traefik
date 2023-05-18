package npool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

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
		EntID            *string
		AppID            string
		UserID           *string
		Path             string
		Method           string
		Arguments        string
		StatusCode       int
		ReqHeaders       string
		RespHeaders      string
		NewValue         string
		Result           *string
		FailReason       string
		ElapsedMillisecs uint32
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

	reqHeaders, _ := json.Marshal(req.Header)
	olq.ReqHeaders = string(reqHeaders)

	type opLogInfo struct {
		EntID string
	}
	type opLogResp struct {
		Info opLogInfo
	}
	var resp *resty.Response

	resp, err = resty.
		New().
		R().
		SetBody(olq).
		SetResult(&opLogResp{}).
		Post(fmt.Sprintf("http://%v/v1/create/oplog", opLogHost))
	if err != nil {
		logger.Warnf("fail create oplog: %v", err)
	}

	req.Body = ioutil.NopCloser(strings.NewReader(string(_body)))
	buffer := bytes.NewBuffer(nil)
	_rw := newMultiWriter(buffer, rw)

	start := time.Now()
	ol.next.ServeHTTP(_rw, req)

	if resp == nil {
		logger.Warnf("fail create oplog: nil resp")
		return
	}

	olr := resp.Result().(*opLogResp)
	if _, err := uuid.Parse(olr.Info.EntID); err != nil {
		logger.Warnf("invalid oplog ent_id %v: %v", resp, err)
		return
	}

	respHeaders, _ := json.Marshal(_rw.Header())
	reqHeaders, _ = json.Marshal(req.Header)
	statusCode := _rw.StatusCode()
	result := "Success"
	if statusCode != http.StatusOK {
		result = "Fail"
	}

	olq = &opLogReq{
		EntID:            &olr.Info.EntID,
		NewValue:         buffer.String(),
		StatusCode:       statusCode,
		ReqHeaders:       string(reqHeaders),
		RespHeaders:      string(respHeaders),
		Result:           &result,
		ElapsedMillisecs: uint32(time.Since(start).Milliseconds()),
	}

	if statusCode != http.StatusOK {
		olq.FailReason = buffer.String()
	}

	resp, err = resty.
		New().
		R().
		SetBody(olq).
		SetResult(&opLogResp{}).
		Post(fmt.Sprintf("http://%v/v1/update/oplog", opLogHost))
	logger.Debugf("update oplog: %v (%v)", resp, err)
}

type multiWriter struct {
	rw         http.ResponseWriter
	multi      io.Writer
	statusCode int
}

func newMultiWriter(buffer *bytes.Buffer, rw http.ResponseWriter) *multiWriter {
	multi := io.MultiWriter(buffer, rw)
	return &multiWriter{
		rw:    rw,
		multi: multi,
	}
}

func (mw *multiWriter) Header() http.Header {
	return mw.rw.Header()
}

func (mw *multiWriter) Write(b []byte) (int, error) {
	return mw.multi.Write(b)
}

func (mw *multiWriter) WriteHeader(statusCode int) {
	mw.statusCode = statusCode
	mw.rw.WriteHeader(statusCode)
}

func (mw *multiWriter) StatusCode() int {
	return mw.statusCode
}
