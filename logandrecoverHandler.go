package main

import (
	"log"
	"net/http"
	"runtime/debug"
)

func logAndRecoverHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from an error in HTTP handler: %v\n", err)
				debug.PrintStack()

				http.Error(w, "Sorry, something went wrong on our end. We're working to fix it!", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
}
