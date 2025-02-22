package handlers

import (
	"fmt"
	"net/http"
)

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				log.Error("error in parse form",
					"error", err)
			}
			log.Info("incoming request",
				"request", fmt.Sprintf("%+v", r))
			h.ServeHTTP(w, r)
		})
}

func RecoveryMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				r := recover()
				if r != nil {
					log.Error("recovering",
						"error", r)
				}
			}()

			h.ServeHTTP(w, r)
		})
}
