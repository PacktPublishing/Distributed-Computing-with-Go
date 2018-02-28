package api

import (
	"encoding/json"
	"net/http"

	"github.com/last-ent/distributed-go/chapter6/goophr/concierge/common"
)

func QueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"code": 405, "msg": "Method Not Allowed."}`))
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var searchTerms []string
	if err := decoder.Decode(&searchTerms); err != nil {
		common.Warn("Unable to parse request." + err.Error())

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"code": 400, "msg": "Unable to parse payload."}`))
		return
	}
}
