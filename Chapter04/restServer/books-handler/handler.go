// restServer/books-handler/handler.go

package booksHandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// StartBooksManager starts a goroutine that changes the state of books (map).
// Primary reason to use a goroutine instead of directly manipulating the books map is to ensure
// that we do not have multiple requests changing books' state simultaneously.
func StartBooksManager(books map[string]bookResource, actionCh <-chan Action) {
	newID := len(books)
	for {
		select {
		case act := <-actionCh:
			switch act.Type {
			case "GET":
				actOnGET(books, act)
			case "POST":
				newID++
				newBookID := fmt.Sprintf("%d", newID)
				actOnPOST(books, act, newBookID)
			case "PUT":
				actOnPUT(books, act)
			case "DELETE":
				actOnDELETE(books, act)
			}
		}
	}
}

/* BookHandler is responsible for ensuring that we process only the valid HTTP Requests.

* GET -> id: Any

* POST -> id: No
*      -> payload: Required

* PUT -> id: Any
*     -> payload: Required

* DELETE -> id: Any
 */
func BookHandler(w http.ResponseWriter, r *http.Request, id string, method string, actionCh chan<- Action) {

	// Ensure that id is set only for valid requests
	isGet := method == "GET"
	idIsSetForPost := method == "POST" && id != ""
	isPutOrPost := method == "PUT" || method == "POST"
	idIsSetForDelPut := (method == "DELETE" || method == "PUT") && id != ""
	if !isGet && !(idIsSetForPost || idIsSetForDelPut || isPutOrPost) {
		writeError(w, http.StatusMethodNotAllowed)
		return
	}

	respCh := make(chan response)
	act := Action{
		Id:      id,
		Type:    method,
		RetChan: respCh,
	}

	// PUT & POST require a properly formed JSON payload
	if isPutOrPost {
		var reqPayload requestPayload
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err := json.Unmarshal(body, &reqPayload); err != nil {
			writeError(w, http.StatusBadRequest)
			return
		}

		act.Payload = reqPayload
	}

	// We have all the data required to process the Request.
	// Time to update the state of books.
	actionCh <- act

	// Wait for respCh to return data after updating the state of books.
	// For all successful Actions, the HTTP status code will either be 200 or 201.
	// Any other status code means that there was an issue with the request.
	var resp response
	if resp = <-respCh; resp.StatusCode > http.StatusCreated {
		writeError(w, resp.StatusCode)
		return
	}

	// We should only log the delete resource and not send it back to user
	if method == "DELETE" {
		log.Println(fmt.Sprintf("Resource ID %s deleted: %+v", id, resp.Books))
		resp = response{
			StatusCode: http.StatusOK,
			Books:      []bookResource{},
		}
	}

	writeResponse(w, resp)
}
