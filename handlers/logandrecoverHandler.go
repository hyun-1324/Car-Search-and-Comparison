package handlers

import (
	"log"
	"net/http"
)

// wraps an HTTP handler to manage panics.
// If a panic occurs, it logs the error then displays a user-friendly message indicating a server issue.
func LogAndRecoverHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from an error in HTTP handler: %v\n", err)
				http.Error(w, "Sorry, something went wrong on our end. We're working to fix it!", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
}
