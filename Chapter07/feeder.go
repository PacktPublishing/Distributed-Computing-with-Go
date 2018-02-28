package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type tPayload struct {
	Token  string `json:"token"`
	Title  string `json:"title"`
	DocID  string `json:"doc_id"`
	LIndex int    `json:"line_index"`
	Index  int    `json:"token_index"`
}

type msgS struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func main() {
	// Searching for "apple" should return Book 1 at the top of search results.
	// Searching for "cake" should return Book 3 at the top.
	for bookX, terms := range map[string][]string{
		"Book 1": []string{"apple", "apple", "cat", "zebra"},
		"Book 2": []string{"banana", "cake", "zebra"},
		"Book 3": []string{"apple", "cake", "cake", "whale"},
	} {
		for lin, term := range terms {
			payload, _ := json.Marshal(tPayload{
				Token:  term,
				Title:  bookX + term,
				DocID:  bookX,
				LIndex: lin,
			})
			resp, err := http.Post(
				"http://localhost:9090/api/index",
				"application/json",
				bytes.NewBuffer(payload),
			)
			if err != nil {
				panic(err)
			}
			body, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()

			var msg msgS
			json.Unmarshal(body, &msg)
			log.Println(msg)
		}
	}
}
