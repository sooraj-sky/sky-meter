package sentry

import (
	"log"

	"github.com/getsentry/sentry-go"
	skyenv "sky-meter/packages/env"
)

// The function SentryInit initializes the Sentry error tracking service.
func SentryInit() {
	allEnv := skyenv.GetEnv()
	mode := allEnv.Mode
	if mode == "dev" {
		sentryEnv := allEnv.SentryDsn
		sentryerr := sentry.Init(sentry.ClientOptions{
			Dsn: sentryEnv,
			// Set TracesSampleRate to 1.0 to capture 100%
			// of transactions for performance monitoring.
			// We recommend adjusting this value in production,
			TracesSampleRate: 1.0,
		})
		if sentryerr != nil {
			log.Fatalf("sentry.Init: %s", sentryerr)
		}

	}
}
