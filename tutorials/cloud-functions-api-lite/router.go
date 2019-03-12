package api

//
// Custom Structured Logger
// ========================
// This example demonstrates how to use middleware.RequestLogger,
// middleware.LogFormatter and middleware.LogEntry to build a structured
// logger using the amazing sirupsen/logrus package as the logging
// backend.
//
// Also: check out https://github.com/pressly/lg for an improved context
// logger with support for HTTP request logging, based on the example
// below.
//

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type ApiRouter struct {
	Handler  *chi.Mux
	logger   *CloudLogger
	reporter *ErrorRecoverer
}

func (ar *ApiRouter) Close() {
	ar.logger.Close()
	ar.reporter.Close()
}

func NewRouter(ctx context.Context) *ApiRouter {

	// Setup the logger backend using sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/sirupsen/logrus

	ar := &ApiRouter{}
	// Routes
	ar.Handler = chi.NewRouter()
	ar.Handler.Use(middleware.RequestID)

	ar.logger = &CloudLogger{}
	// TODO get logname from context?
	ar.Handler.Use(ar.logger.NewStructuredLogger("api-log"))

	reporter := NewErrorRecoverer(ctx)
	ar.Handler.Use(reporter.Middleware)

	ar.Handler.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	ar.Handler.Get("/wait", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		LogEntrySetField(r, "wait", true)
		w.Write([]byte("hi"))
	})
	ar.Handler.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("oops")
	})
	return ar
}
