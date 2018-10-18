package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// PurgeParam :
type PurgeParam struct {
	PurgeCacheURL string `json:"purgeCacheURL"`
}

// HandlePurgeRequest :
func HandlePurgeRequest(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		log.Printf("Body is nil : %v", r.Body)
		w.WriteHeader(http.StatusOK)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if len(b) == 0 {
		log.Print("Request body is empty.")
		w.WriteHeader(http.StatusOK)
		return
	}
	if err != nil {
		log.Printf("Failed to read request body: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}

	// unmarshal
	var param PurgeParam
	err = json.Unmarshal(b, &param)
	if err != nil {
		log.Printf("Failed to unmarshal params: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}

	// send purge request
	req, err := http.NewRequest("PURGE", param.PurgeCacheURL, nil)
	if err != nil {
		log.Printf("Failed to init request: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send purge request: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(http.StatusOK)
	return
}
