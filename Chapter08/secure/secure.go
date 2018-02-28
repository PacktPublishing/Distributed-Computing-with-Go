package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

var authTokens = map[string]string{
	"AUTH-TOKEN-1": "User 1",
	"AUTH-TOKEN-2": "User 2",
}

// getAuthorizedUser tries to retrieve user for the given token.
func getAuthorizedUser(token string) (string, error) {
	var err error

	user, valid := authTokens[token]
	if !valid {
		err = fmt.Errorf("Auth token '%s' does not exist.", token)
	}

	return user, err
}

// isAuthorized checks request to ensure that it has Authorization header
// with defined value: "Bearer AUTH-TOKEN"
func isAuthorized(r *http.Request) bool {
	rawToken := r.Header["Authorization"]
	if len(rawToken) != 1 {
		return false
	}

	authToken := strings.Split(rawToken[0], " ")
	if !(len(authToken) == 2 && authToken[0] == "Bearer") {
		return false
	}

	user, err := getAuthorizedUser(authToken[1])
	if err != nil {
		log.Printf("Error: %s", err)
		return false
	}

	log.Printf("Successful request made by '%s'", user)
	return true
}

var success = []byte("Received authorized request.")
var failure = []byte("Received unauthorized request.")

func requestHandler(w http.ResponseWriter, r *http.Request) {
	if isAuthorized(r) {
		w.Write(success)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(failure)
	}
}

func main() {
	http.HandleFunc("/", requestHandler)
	fmt.Println("Starting server @ http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
