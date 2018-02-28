// restServer/books-handler/actions.go

package booksHandler

import (
	"net/http"
)

// actOn{GET, POST, DELETE, PUT} functions return Response based on specific Request type.

func actOnGET(books map[string]bookResource, act Action) {
	// These initialized values cover the case:
	// Request asked for an id that doesn't exist.
	status := http.StatusNotFound
	bookResult := []bookResource{}

	if act.Id == "" {

		// Request asked for all books.
		status = http.StatusOK
		for _, book := range books {
			bookResult = append(bookResult, book)
		}
	} else if book, exists := books[act.Id]; exists {

		// Request asked for a specific book and the id exists.
		status = http.StatusOK
		bookResult = []bookResource{book}
	}

	act.RetChan <- response{
		StatusCode: status,
		Books:      bookResult,
	}
}

func actOnDELETE(books map[string]bookResource, act Action) {
	book, exists := books[act.Id]
	delete(books, act.Id)

	if !exists {
		book = bookResource{}
	}

	// Return the deleted book if it exists else return an empty book.
	act.RetChan <- response{
		StatusCode: http.StatusOK,
		Books:      []bookResource{book},
	}
}

func actOnPUT(books map[string]bookResource, act Action) {
	// These initialized values cover the case:
	// Request asked for an id that doesn't exist.
	status := http.StatusNotFound
	bookResult := []bookResource{}

	// If the id exists, update its values with the values from the payload.
	if book, exists := books[act.Id]; exists {
		book.Link = act.Payload.Link
		book.Title = act.Payload.Title
		books[act.Id] = book

		status = http.StatusOK
		bookResult = []bookResource{books[act.Id]}
	}

	// Return status and updated resource.
	act.RetChan <- response{
		StatusCode: status,
		Books:      bookResult,
	}

}

func actOnPOST(books map[string]bookResource, act Action, newID string) {
	// Add the new book to `books`.
	books[newID] = bookResource{
		Id:    newID,
		Link:  act.Payload.Link,
		Title: act.Payload.Title,
	}

	act.RetChan <- response{
		StatusCode: http.StatusCreated,
		Books:      []bookResource{books[newID]},
	}
}
