// restClient.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type bookResource struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

func main() {
	// GET
	fmt.Println("Making GET call.")
	// It is possible that we might have error while making an HTTP request
	// due to too many redirects or HTTP protocol error. We should check for this eventuality.
	resp, err := http.Get("http://localhost:8080/api/books")
	if err != nil {
		fmt.Println("Error while making GET call.", err)
		return
	}

	fmt.Printf("%+v\n\n", resp)

	// The response body is a data stream from the server we got the response back from.
	// This data stream is not in a useable format yet.
	// We need to read it from the server and convert it into a byte stream.
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var books []bookResource
	json.Unmarshal(body, &books)

	fmt.Println(books)
	fmt.Println("\n")

	// POST
	payload, _ := json.Marshal(bookResource{
		Title: "New Book",
		Link:  "http://new-book.com",
	})

	fmt.Println("Making POST call.")
	resp, err = http.Post(
		"http://localhost:8080/api/books/",
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n\n", resp)

	body, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var book bookResource
	json.Unmarshal(body, &book)

	fmt.Println(book)

	fmt.Println("\n")
}
