package main

import (
	"log"
	"net/http"
)

func reqHandler(w http.ResponseWriter, r *http.Request) {
	books := map[string]string{
		"book1": `apple apple cat zebra`,
		"book2": `banana cake zebra`,
		"book3": `apple cake cake whale`,
	}

	bookID := r.URL.Path[1:]
	book, _ := books[bookID]
	w.Write([]byte(book))
}

func main() {

	log.Println("Starting File Server on Port :9876...")
	http.HandleFunc("/", reqHandler)
	http.ListenAndServe(":9876", nil)
}
