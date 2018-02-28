package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const passMark = "\u2713"
const failMark = "\u2717"

func assertResponseEqual(t *testing.T, expected string, actual string) {
	t.Helper() // comment this line to see tests fail due to `if expected != actual`
	if expected != actual {
		t.Errorf("%s != %s %s", expected, actual, failMark)
	} else {
		t.Logf("%s == %s %s", expected, actual, passMark)
	}
}

func TestServer(t *testing.T) {
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				path := r.RequestURI
				if path == "/1" {
					w.Write([]byte("Got 1."))
				} else {
					w.Write([]byte("Got None."))
				}
			}))
	defer testServer.Close()

	for _, testCase := range []struct {
		Name     string
		Path     string
		Expected string
	}{
		{"Request correct URL", "/1", "Got 1."},
		{"Request incorrect URL", "/12345", "Got None."},
	} {
		t.Run(testCase.Name, func(t *testing.T) {
			res, err := http.Get(testServer.URL + testCase.Path)
			if err != nil {
				t.Fatal(err)
			}

			actual, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			assertResponseEqual(t, testCase.Expected, fmt.Sprintf("%s", actual))
		})
	}
	t.Run("Fail for no reason", func(t *testing.T) {
		assertResponseEqual(t, "+", "-")
	})
}
