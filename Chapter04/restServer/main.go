// restServer/main.go

package main

import (
	"fmt"
	"log"
	"net/http"

	booksHandler "github.com/last-ent/distributed-go/chapter4/restServer/books-handler"
)

func main() {
	// Get state (map) for books available on REST server.
	books := booksHandler.GetBooks()
	log.Println(fmt.Sprintf("%+v", books))

	actionCh := make(chan booksHandler.Action)

	// Start goroutine responsible for handling interaction with the books map
	go booksHandler.StartBooksManager(books, actionCh)

	http.HandleFunc("/api/books/", booksHandler.MakeHandler(booksHandler.BookHandler, "/api/books/", actionCh))

	log.Println("Starting server at port 8080...")
	http.ListenAndServe(":8080", nil)
}
