package rest

import (
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"log"
	"net/http"
	"time"
)

func SentryMiddleware(dsn, env string, debug bool, release string) func(next http.Handler) http.Handler {
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
		Environment:      env,
		Release:          release,
		Debug:            debug,
	})
	if err != nil {
		log.Fatal(err)
	}
	hub := sentry.NewHub(client, sentry.NewScope())
	sentryHandler := sentryhttp.New(sentryhttp.Options{Repanic: true, Timeout: time.Minute, WaitForDelivery: true})
	handler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(sentry.SetHubOnContext(r.Context(), hub))
			sentryHandler.Handle(next).ServeHTTP(w, r)
		})
	}
	return handler
}
