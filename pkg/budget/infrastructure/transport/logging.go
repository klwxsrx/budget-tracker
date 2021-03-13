package transport

import (
	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
	"io/ioutil"
	"net/http"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	code int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func getLoggingMiddleware(l logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lrw := newLoggingResponseWriter(w)
			next.ServeHTTP(lrw, r)

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				body = nil
			}
			loggerWithFields := l.With(logger.Fields{
				"method":       r.Method,
				"url":          r.RequestURI,
				"body":         string(body),
				"responseCode": lrw.code,
			})
			if lrw.code == 500 {
				loggerWithFields.Error("internal server error")
			} else {
				loggerWithFields.Info("request handled")
			}
		})
	}
}
