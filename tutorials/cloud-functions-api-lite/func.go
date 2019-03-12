package api

import (
	"context"
	"net/http"
)

var router *ApiRouter

func init() {
	ctx := context.Background()
	router = NewRouter(ctx)
}

// API acts as a single cloud function entry point to manage a limited API surface
func API(w http.ResponseWriter, r *http.Request) {
	router.Handler.ServeHTTP(w, r)
}
