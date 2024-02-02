package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// prevent XSS with a strict content security policy
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		// recommended when setting CSP
		w.Header().Set("X-XSS-Protection", "0")

		// hide URL path or query values in non same origin requests
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")

		// no content sniffing... looking at you internet explorer
		w.Header().Set("X-Content-Type-Options", "no-sniff")

		// no clickjacking
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

// implement as method to get access to infoLog
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// deferred function will always run when panic occurs, as go
		// will unwind the goroutine serving the request
		defer func() {
			// retrieve the error that caused the panic
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err)) // ensure error type is string
			}
		}()
		next.ServeHTTP(w, r)
	})
}
