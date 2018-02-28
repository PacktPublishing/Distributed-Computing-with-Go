package api

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/last-ent/distributed-go/chapter8/goophr/concierge/common"
)

type payload struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type document struct {
	Doc   string `json:"-"`
	Title string `json:"title"`
	DocID string `json:"-"`
	URL string `json:"url"`
}

type token struct {
	Line   string `json:"-"`
	Token  string `json:"token"`
	Title  string `json:"title"`
	DocID  string `json:"doc_id"`
	LIndex int    `json:"line_index"`
	Index  int    `json:"token_index"`
}

type dMsg struct {
	DocID string
	Ch    chan document
}

type lMsg struct {
	LIndex int
	DocID  string
	Ch     chan string
}

type lMeta struct {
	LIndex int
	DocID  string
	Line   string
}

type dAllMsg struct {
	Ch chan []document
}

// done signals all listening goroutines to stop.
var done chan bool

// dGetCh is used to retrieve a single document from store.
var dGetCh chan dMsg

// lGetCh is used to retrieve a single line from store.
var lGetCh chan lMsg

// lStoreCh is used to put a line into store.
var lStoreCh chan lMeta

// iAddCh is used to add token to index (Librarian).
var iAddCh chan token

// dStoreCh is used to put a document into store.
var dStoreCh chan document

// dProcessCh is used to process a document and convert it to tokens.
var dProcessCh chan document

// dGetAllCh is used to retrieve all documents in store.
var dGetAllCh chan dAllMsg

// pProcessCh is used to process the /feeder's payload and start the indexing process.
var pProcessCh chan payload

// StartFeederSystem initializes all channels and starts all goroutines.
// We are using a standard function instead of `init()`
// because we don't want the channels & goroutines to be initialized during testing.
// Unless explicitly required by a particular test.
func StartFeederSystem() {
	done = make(chan bool)

	dGetCh = make(chan dMsg, 8)
	dGetAllCh = make(chan dAllMsg)

	iAddCh = make(chan token, 8)
	pProcessCh = make(chan payload, 8)

	dStoreCh = make(chan document, 8)
	dProcessCh = make(chan document, 8)
	lGetCh = make(chan lMsg)
	lStoreCh = make(chan lMeta, 8)

	for i := 0; i < 4; i++ {
		go indexAdder(iAddCh, done)
		go docProcessor(pProcessCh, dStoreCh, dProcessCh, done)
		go indexProcessor(dProcessCh, lStoreCh, iAddCh, done)
	}

	go docStore(dStoreCh, dGetCh, dGetAllCh, done)
	go lineStore(lStoreCh, lGetCh, done)
}

func getTokIndexer(tok string) string {
	t := tok[0]
	indexerURL := librarianEndpoints["*"]

	if t >= 'a' && t <= 'm' {
		indexerURL = librarianEndpoints["a-m"]
	} else if t >= 'n' && t <= 'z' {
		indexerURL = librarianEndpoints["n-z"]
	}

	return indexerURL + "/index"
}

// indexAdder adds token to index (Librarian).
func indexAdder(ch chan token, done chan bool) {
	for {
		select {
		case tok := <-ch:
			fmt.Println("adding to librarian:", tok.Token)
			serializedPayload, err := json.Marshal(tok)
			if err != nil {
				common.Warn("Unable to serialize the payload. Error:" + err.Error())
			}

			indexerURL := getTokIndexer(tok.Token)
			resp, err := http.Post(indexerURL, "application/json", bytes.NewBuffer(serializedPayload))
			if err != nil {
				common.Warn("Error while posting to Librarian. Error:" + err.Error())
				continue
			}

			rawMsg, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				common.Warn("Error while reading response body:" + err.Error())
				continue
			}

			msg := string(rawMsg)
			if resp.StatusCode > 200 {
				common.Warn("Error while posting to Librarian. Msg:" + msg)
			} else {
				common.Log("Request was posted to Librairan. Msg:" + msg)
			}

		case <-done:
			common.Log("Exiting indexAdder.")
			return
		}
	}
}

// lineStore maintains a catalog of all lines for all documents being indexed.
func lineStore(ch chan lMeta, callback chan lMsg, done chan bool) {
	store := map[string]string{}
	for {
		select {
		case line := <-ch:
			id := fmt.Sprintf("%s-%d", line.DocID, line.LIndex)
			store[id] = line.Line

		case ch := <-callback:
			line := ""
			id := fmt.Sprintf("%s-%d", ch.DocID, ch.LIndex)
			if l, exists := store[id]; exists {
				line = l
			}
			ch.Ch <- line
		case <-done:
			common.Log("Exiting docStore.")
			return
		}
	}
}

// indexProcessor is responsible for converting a document into tokens for indexing.
func indexProcessor(ch chan document, lStoreCh chan lMeta, iAddCh chan token, done chan bool) {
	for {
		select {
		case doc := <-ch:
			docLines := strings.Split(doc.Doc, "\n")

			lin := 0
			for _, line := range docLines {
				if strings.TrimSpace(line) == "" {
					continue
				}

				lStoreCh <- lMeta{
					LIndex: lin,
					Line:   line,
					DocID:  doc.DocID,
				}

				index := 0
				words := strings.Fields(line)
				for _, word := range words {
					if tok, valid := common.SimplifyToken(word); valid {
						iAddCh <- token{
							Token:  tok,
							LIndex: lin,
							Line:   line,
							Index:  index,
							DocID:  doc.DocID,
							Title:  doc.Title,
						}
						index++
					}
				}
				lin++
			}

		case <-done:
			common.Log("Exiting indexProcessor.")
			return
		}
	}
}

// docStore maintains a catalog of all documents being indexed.
func docStore(add chan document, get chan dMsg, dGetAllCh chan dAllMsg, done chan bool) {
	store := map[string]document{}

	for {
		select {
		case doc := <-add:
			store[doc.DocID] = doc
		case m := <-get:
			m.Ch <- store[m.DocID]
		case ch := <-dGetAllCh:
			docs := []document{}
			for _, doc := range store {
				docs = append(docs, doc)
			}
			ch.Ch <- docs
		case <-done:
			common.Log("Exiting docStore.")
			return
		}
	}
}

// docProcessor processes new document payloads.
func docProcessor(in chan payload, dStoreCh chan document, dProcessCh chan document, done chan bool) {
	for {
		select {
		case newDoc := <-in:
			var err error
			doc := ""

			if doc, err = getFile(newDoc.URL); err != nil {
				common.Warn(err.Error())
				continue
			}

			titleID := getTitleHash(newDoc.Title)
			msg := document{
				Doc:   doc,
				DocID: titleID,
				URL: newDoc.URL,
				Title: newDoc.Title,
			}

			dStoreCh <- msg
			dProcessCh <- msg
		case <-done:
			common.Log("Exiting docProcessor.")
			return
		}
	}
}

// getTitleHash returns a new hash ID everytime it is called.
// Based on: https://gobyexample.com/sha1-hashes
func getTitleHash(title string) string {
	hash := sha1.New()
	title = strings.ToLower(title)

	str := fmt.Sprintf("%s-%s", time.Now(), title)
	hash.Write([]byte(str))

	hByte := hash.Sum(nil)

	return fmt.Sprintf("%x", hByte)
}

// getFile returns file content after retrieving it from URL.
func getFile(URL string) (string, error) {
	var res *http.Response
	var err error

	if res, err = http.Get(URL); err != nil {
		errMsg := fmt.Errorf("Unable to retrieve URL: %s.\nError: %s", URL, err)

		return "", errMsg

	}
	if res.StatusCode > 200 {
		errMsg := fmt.Errorf("Unable to retrieve URL: %s.\nStatus Code: %d", URL, res.StatusCode)

		return "", errMsg
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		errMsg := fmt.Errorf("Error while reading response: URL: %s.\nError: %s", URL, res.StatusCode, err.Error())

		return "", errMsg
	}

	return string(body), nil
}

// FeedHandler start processing the payload which contains the file to index.
func FeedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ch := make(chan []document)
		dGetAllCh <- dAllMsg{Ch: ch}
		docs := <-ch
		close(ch)

		if serializedPayload, err := json.Marshal(docs); err == nil {
			w.Write(serializedPayload)
		} else {
			common.Warn("Unable to serialize all docs: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"code": 500, "msg": "Error occurred while trying to retrieve documents."}`))
		}
		return
	} else if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"code": 405, "msg": "Method Not Allowed."}`))
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var newDoc payload
	decoder.Decode(&newDoc)
	pProcessCh <- newDoc

	w.Write([]byte(`{"code": 200, "msg": "Request is being processed."}`))
}
