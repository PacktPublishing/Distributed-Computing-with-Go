package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/last-ent/distributed-go/chapter8/goophr/librarian/api"
	"github.com/last-ent/distributed-go/chapter8/goophr/librarian/common"
)

func main() {
	common.Log("Adding API handlers...")
	http.HandleFunc("/api/index", api.IndexHandler)
	http.HandleFunc("/api/query", api.QueryHandler)

	common.Log("Starting index...")
	api.StartIndexSystem()

	port := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	common.Log(fmt.Sprintf("Starting Goophr Librarian server on port %s...", port))
	http.ListenAndServe(port, nil)
}
