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
package main

import (
	"context"
	"net/http"

	"api"
)

func main() {
	ctx := context.Background()
	router := api.NewRouter(ctx)
	defer router.Close()
	http.ListenAndServe(":8090", router.Handler)
}
