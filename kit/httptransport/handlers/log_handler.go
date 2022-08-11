package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/iotexproject/Bumblebee/conf/log"
	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"
	"github.com/iotexproject/Bumblebee/kit/metax"
	"github.com/iotexproject/Bumblebee/x/misc/timer"
)

func LogHandler() func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return &loggerHandler{
			next: handler,
		}
	}
}

type loggerHandler struct {
	next http.Handler
}

type LoggerResponseWriter struct {
	rw         http.ResponseWriter
	written    bool
	statusCode int
	err        error
}

func (rw *LoggerResponseWriter) Header() http.Header { return rw.rw.Header() }

func (rw *LoggerResponseWriter) WriteErr(err error) { rw.err = err }

func (rw *LoggerResponseWriter) WriteHeader(sc int) {
	if rw.written {
		return
	}
	rw.rw.WriteHeader(sc)
	rw.statusCode = sc
	rw.written = true
}

func (rw *LoggerResponseWriter) Write(data []byte) (int, error) {
	if rw.err != nil && rw.statusCode >= http.StatusBadRequest {
		rw.err = errors.New(string(data))
	}
	return rw.rw.Write(data)
}

func (h *loggerHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	cost := timer.Start()
	reqID := req.Header.Get(httpx.HeaderRequestID)
	if reqID == "" {
		reqID = uuid.New().String()
	}

	var (
		writer   = &LoggerResponseWriter{rw: rw}
		logger   = log.FromContext(req.Context())
		level, _ = log.ParseLevel(strings.ToLower(req.Header.Get("x-log-level")))
	)

	if level == log.PanicLevel {
		level = log.TraceLevel
	}

	defer func() {
		header := req.Header
		fields := []interface{}{
			"tag", "access",
			"cost", fmt.Sprintf("%0.3fms", float64(cost()/time.Millisecond)),
			"remote_ip", httpx.ClientIP(req),
			"method", req.Method[0:3],
			"request_url", req.URL.String(),
			"user_agent", header.Get(httpx.HeaderUserAgent),
			"status", writer.statusCode,
		}
		if writer.err != nil {
			if writer.statusCode >= http.StatusInternalServerError {
				if level >= log.ErrorLevel {
					logger.WithValues(fields).Error(writer.err)
				}
			} else {
				if level >= log.WarnLevel {
					logger.WithValues(fields).Warn(writer.err)
				}
			}
		} else {
			if level >= log.InfoLevel {
				logger.WithValues(fields).Info("")
			}
		}
	}()

	h.next.ServeHTTP(
		writer,
		req.WithContext(
			metax.ContextWithMeta(req.Context(), metax.ParseMeta(reqID)),
		),
	)
}
