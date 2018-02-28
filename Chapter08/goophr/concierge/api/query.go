package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/last-ent/distributed-go/chapter8/goophr/concierge/common"
)

var librarianEndpoints = map[string]string{}

func init() {
	librarianEndpoints["a-m"] = os.Getenv("LIB_A_M")
	librarianEndpoints["n-z"] = os.Getenv("LIB_N_Z")
	librarianEndpoints["*"] = os.Getenv("LIB_OTHERS")
}

type docs struct {
	DocID string `json:"doc_id"`
	Score int    `json:"doc_score"`
}

type queryResult struct {
	Count int    `json:"count"`
	Data  []docs `json:"data"`
}

func queryLibrarian(endpoint string, stBytes io.Reader, ch chan<- queryResult) {
	resp, err := http.Post(
		endpoint+"/query",
		"application/json",
		stBytes,
	)
	if err != nil {
		common.Warn(fmt.Sprintf("%s -> %+v", endpoint, err))
		ch <- queryResult{}
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var qr queryResult
	json.Unmarshal(body, &qr)
	log.Println(fmt.Sprintf("%s -> %#v", endpoint, qr))
	ch <- qr
}

func getResultsMap(ch <-chan queryResult) map[string]int {
	results := []docs{}
	for range librarianEndpoints {
		if result := <-ch; result.Count > 0 {
			results = append(results, result.Data...)
		}
	}

	resultsMap := map[string]int{}
	for _, doc := range results {
			docID := doc.DocID
			score := doc.Score
			if _, exists := resultsMap[docID]; !exists {
				resultsMap[docID] = 0
			}
			resultsMap[docID] = resultsMap[docID] + score
		}

	return resultsMap
}

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

	st, err := json.Marshal(searchTerms)
	if err != nil {
		panic(err)
	}
	stBytes := bytes.NewBuffer(st)

	resultsCh := make(chan queryResult)

	for _, le := range librarianEndpoints {
		func(endpoint string) {
			go queryLibrarian(endpoint, stBytes, resultsCh)
		}(le)
	}

	resultsMap := getResultsMap(resultsCh)
	close(resultsCh)

	sortedResults := sortResults(resultsMap)

	payload, _ := json.Marshal(sortedResults)
	w.Header().Add("Content-Type", "application/json")
	w.Write(payload)

	fmt.Printf("%#v\n", sortedResults))
}

func sortResults(rm map[string]int) []document {
	scoreMap := map[int][]document{}
	ch := make(chan document)
	
	for docID, score := range rm {
		if _, exists := scoreMap[score]; !exists {
			scoreMap[score] = []document{}
		}

		dGetCh <- dMsg{
			DocID: docID,
			Ch:    ch,
		}
		doc := <-ch

		scoreMap[score] = append(scoreMap[score], doc)
	}

	close(ch)

	scores := []int{}
	for score := range scoreMap {
		scores = append(scores, score)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(scores)))

	sortedResults := []document{}
	for _, score := range scores {
		resDocs := scoreMap[score]
		sortedResults = append(sortedResults, resDocs...)
	}

	return sortedResults
}
