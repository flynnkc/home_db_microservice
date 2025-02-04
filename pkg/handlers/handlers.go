package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
)

var log *slog.Logger = slog.Default()

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Info("incoming request",
				"request", fmt.Sprintf("%+v", r))
			h.ServeHTTP(w, r)
			log.Info("after request")
		})
}

func TestMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Test")
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

func MysqlReadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error("error parsing form",
			"error", err.Error(),
			"request", *r)
	}

	if err != nil {
		log.Error("error reading request",
			"error", err.Error(),
			"request", *r)
	}

	fmt.Println(r.FormValue("foo"))
	w.WriteHeader(200)
	w.Write([]byte("success\n"))
}
