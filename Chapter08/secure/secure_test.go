package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsAuthorizedSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Error("Unable to create request")
	}

	req.Header["Authorization"] = []string{"Bearer AUTH-TOKEN-1"}

	if isAuthorized(req) {
		t.Log("Request with correct Auth token was correctly processed.")
	} else {
		t.Error("Request with correct Auth token failed.")
	}
}

func TestIsAuthorizedFailTokenType(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Error("Unable to create request")
	}

	req.Header["Authorization"] = []string{"Token AUTH-TOKEN-1"}

	if isAuthorized(req) {
		t.Error("Request with incorrect Auth token type was successfully processed.")
	} else {
		t.Log("Request with incorrect Auth token type failed as expected.")
	}
}

func TestIsAuthorizedFailToken(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Error("Unable to create request")
	}

	req.Header["Authorization"] = []string{"Token WRONG-AUTH-TOKEN"}

	if isAuthorized(req) {
		t.Error("Request with incorrect Auth token was successfully processed.")
	} else {
		t.Log("Request with incorrect Auth token failed as expected.")
	}
}

func TestRequestHandlerFailToken(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Error("Unable to create request")
	}

	req.Header["Authorization"] = []string{"Token WRONG-AUTH-TOKEN"}

	// http.ResponseWriter it is an interface hence we use
	// httptest.NewRecorder which implements the interface http.ResponseWriter
	rr := httptest.NewRecorder()
	requestHandler(rr, req)

	if rr.Code == 401 {
		t.Log("Request with incorrect Auth token failed as expected.")
	} else {
		t.Error("Request with incorrect Auth token was successfully processed.")
	}
}

func TestGetAuthorizedUser(t *testing.T) {
	if user, err := getAuthorizedUser("AUTH-TOKEN-2"); err != nil {
		t.Errorf("Couldn't find User 2. Error: %s", err)
	} else if user != "User 2" {
		t.Errorf("Found incorrect user: %s", user)
	} else {
		t.Log("Found User 2.")
	}
}

func TestGetAuthorizedUserFail(t *testing.T) {
	if user, err := getAuthorizedUser("WRONG-AUTH-TOKEN"); err == nil {
		t.Errorf("Found user for invalid token!. User: %s", user)
	} else if err.Error() != "Auth token 'WRONG-AUTH-TOKEN' does not exist." {
		t.Errorf("Error message does not match.")
	} else {
		t.Log("Got expected error message for invalid auth token")
	}
}
