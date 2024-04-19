package main

import (
	"log"
	"net/http"

	"cars/handlers"
)

func main() {
	http.HandleFunc("/", handlers.LogAndRecoverHandler(handlers.MainHandler))
	http.HandleFunc("/result", handlers.LogAndRecoverHandler(handlers.ResultHandler))
	http.HandleFunc("/comparison", handlers.LogAndRecoverHandler(handlers.ComparisonHandler))
	http.HandleFunc("/popup", handlers.LogAndRecoverHandler(handlers.PopupHandler))
	http.HandleFunc("/download_txt", handlers.LogAndRecoverHandler(handlers.DownloadTxt))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Staring server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
