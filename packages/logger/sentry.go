package sentry

import (
	"log"
	"os"

	"github.com/getsentry/sentry-go"
)

func SentryInit() {
	mode := os.Getenv("mode")
	if mode == "dev" {
		sentenv := os.Getenv("sentry_dsn")
		if sentenv == "" {
			log.Fatal("Please specify the sentry_dsn as environment variable, e.g. env sentry_dsn=https://your-dentry-dsn.com go run cmd/main.go")
		}
		senterr := sentry.Init(sentry.ClientOptions{
			Dsn: sentenv,
			// Set TracesSampleRate to 1.0 to capture 100%
			// of transactions for performance monitoring.
			// We recommend adjusting this value in production,
			TracesSampleRate: 1.0,
		})
		if senterr != nil {
			log.Fatalf("sentry.Init: %s", senterr)
		}

	}
}
