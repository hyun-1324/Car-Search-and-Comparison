package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", logAndRecoverHandler(mainHandler))
	http.HandleFunc("/result", logAndRecoverHandler(resultHandler))
	http.HandleFunc("/comparison", logAndRecoverHandler(comparisonHandler))
	http.HandleFunc("/popup", logAndRecoverHandler(popupHandler))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Staring server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
