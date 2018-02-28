package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTitleHash(t *testing.T) {

	h1 := getTitleHash("A-Title")
	h2 := getTitleHash("Diff Title")
	hDup := getTitleHash("A-Title")

	for _, tc := range []struct {
		name     string
		hashes   []string
		expected bool
	}{
		{"Different Titles", []string{h1, h2}, false},
		{"Duplicate Titles", []string{h1, hDup}, false},
		{"Same hashes", []string{h2, h2}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.hashes[0] == tc.hashes[1]
			if actual != tc.expected {
				t.Error(actual, tc.expected, tc.hashes)
			}
		})
	}
}

func TestGetFile(t *testing.T) {
	doc := "Server returned text!"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(doc))
	}))
	defer testServer.Close()

	rDoc, err := getFile(testServer.URL)
	if err != nil {
		t.Error("Error while retrieving document", err)
	}
	if doc != rDoc {
		t.Error(doc, "!=", rDoc)
	}
}

func TestIndexProcessor(t *testing.T) {
	ch1 := make(chan document, 1)
	ch2 := make(chan lMeta, 1)
	ch3 := make(chan token, 3)
	done := make(chan bool)

	go indexProcessor(ch1, ch2, ch3, done)

	ch1 <- document{
		DocID: "a-hash",
		Title: "a-title",
		Doc:   "Golang Programming rocks!",
	}

	for i, tc := range []string{
		"golang", "programming", "rocks",
	} {
		t.Run(fmt.Sprintf("Testing if '%s' is returned. at index: %d", tc, i+1), func(t *testing.T) {
			tok := <-ch3
			if tok.Token != tc {
				t.Error(tok.Token, "!=", tc)
			}
			if tok.Index != i+1 {
				t.Error(tok.Index, "!=", i+1)
			}
		})
	}
	close(done)

}
