package api

// The original work was derived from Goji's middleware, source:
// https://github.com/zenazn/goji/tree/master/web/middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"cloud.google.com/go/errorreporting"
	"github.com/sirupsen/logrus"
)

type ErrorRecoverer struct {
	errorClient *errorreporting.Client
}

func (e *ErrorRecoverer) Close() {
	if (e.errorClient != &errorreporting.Client{}) {
		e.errorClient.Close()
	}
}

func NewErrorRecoverer(ctx context.Context) (e *ErrorRecoverer) {
	e = &ErrorRecoverer{}
	projectID := getProject()
	var err error
	e.errorClient, err = errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: "myservice",
		OnError: func(err error) {
			log.Printf("Could not log error: %v", err)
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	return
}

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a request ID if one is provided.
//
// Alternatively, look at https://github.com/pressly/lg middleware pkgs.
func (er *ErrorRecoverer) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {

				logEntry := GetLogEntry(r)
				if logEntry != nil {
					fmt.Println("********************************* have log")
					// logEntry.Panic(rvr, debug.Stack())
					// calling a log-entry with panic level causes a panic within a panic
					logEntry.WithFields(logrus.Fields{
						"stack": string(debug.Stack()),
						"panic": fmt.Sprintf("%+v", rvr),
					}).Error(rvr)
				} else {
					fmt.Println("-------------------------------------- no log")
					fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
					debug.PrintStack()
				}
				fmt.Println("#############   calling error client")
				er.errorClient.Report(errorreporting.Entry{
					Error: fmt.Errorf("Panic: %v+", rvr),
					Req:   r,
					Stack: debug.Stack(),
				})
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
