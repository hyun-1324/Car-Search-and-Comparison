package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/result", resultHandler)
	http.HandleFunc("/comparison", comparisonHandler)
	http.HandleFunc("/popup", popupHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
