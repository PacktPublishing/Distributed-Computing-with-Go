// restServer/books-handler/common.go

package booksHandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// bookResource is used to hold all data needed to represent a Book resource in the books map.
type bookResource struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

// requestPayload is used to parse request's Payload. We ignore Id field for simplicity.
type requestPayload struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

// response struct consists of all the information required to create the correct HTTP response.
type response struct {
	StatusCode int
	Books      []bookResource
}

// Action struct is used to send data to the goroutine managing the state (map) of books.
// RetChan allows us to send data back to the Handler function so that we can complete the HTTP request.
type Action struct {
	Id      string
	Type    string
	Payload requestPayload
	RetChan chan<- response
}

// GetBooks is used to get the initial state of books represented by a map.
func GetBooks() map[string]bookResource {
	books := map[string]bookResource{}
	for i := 1; i < 6; i++ {
		id := fmt.Sprintf("%d", i)
		books[id] = bookResource{
			Id:    id,
			Title: fmt.Sprintf("Book-%s", id),
			Link:  fmt.Sprintf("http://link-to-book%s.com", id),
		}
	}
	return books
}

// MakeHandler shows a common pattern used reduce duplicated code.
func MakeHandler(fn func(http.ResponseWriter, *http.Request, string, string, chan<- Action),
	endpoint string, actionCh chan<- Action) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		msg := fmt.Sprintf("Received request [%s] for path: [%s]", method, path)
		log.Println(msg)

		id := path[len(endpoint):]
		log.Println("ID is ", id)
		fn(w, r, id, method, actionCh)
	}
}

// writeResponse uses the pattern similar to MakeHandler.
func writeResponse(w http.ResponseWriter, resp response) {
	var err error
	var serializedPayload []byte

	if len(resp.Books) == 1 {
		serializedPayload, err = json.Marshal(resp.Books[0])
	} else {
		serializedPayload, err = json.Marshal(resp.Books)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError)
		fmt.Println("Error while serializing payload: ", err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(serializedPayload)
	}
}

// writeError allows us to return error message in JSON format.
func writeError(w http.ResponseWriter, statusCode int) {
	jsonMsg := struct {
		Msg  string `json:"msg"`
		Code int    `json:"code"`
	}{
		Code: statusCode,
		Msg:  http.StatusText(statusCode),
	}

	if serializedPayload, err := json.Marshal(jsonMsg); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Error while serializing payload: ", err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(serializedPayload)
	}
}
