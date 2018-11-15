package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"
)

// ConvertParam :
type ConvertParam struct {
	Username string `json:"username"`
	GraphID  string `json:"graphID"`
	Date     string `json:"date"`
	Mode     string `json:"mode"`
	Stage    string `json:"stage"`
	Hash     string `json:"hash"`
}

// HandleSVGConvert :
func HandleSVGConvert(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Printf("Failed to unmarshal params: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := u.Query()

	var username string
	if query["username"] != nil {
		username = query["username"][0]
	}
	var graphID string
	if query["graphID"] != nil {
		graphID = query["graphID"][0]
	}
	var date string
	if query["date"] != nil {
		date = query["date"][0]
	}
	var mode string
	if query["mode"] != nil {
		mode = query["mode"][0]
	}
	var stage string
	if query["stage"] != nil {
		stage = query["stage"][0]
	}
	var hash string
	if query["hash"] != nil {
		hash = query["hash"][0]
	}

	param := &ConvertParam{
		Username: username,
		GraphID:  graphID,
		Date:     date,
		Mode:     mode,
		Stage:    stage,
		Hash:     hash,
	}

	// treat params
	paramsStr := "?"
	if param.Date != "" {
		paramsStr = fmt.Sprintf("%sdate=%s", paramsStr, param.Date)
	}
	if param.Mode != "" {
		paramsStr = fmt.Sprintf("%s&mode=%s", paramsStr, param.Mode)
	}

	// get svg
	var url string
	if param.Stage == "dev" {
		url = fmt.Sprintf("https://www-dev-215102.appspot.com/v1/users/%s/graphs/%s%s", param.Username, param.GraphID, paramsStr)
	} else {
		url = fmt.Sprintf("https://pixe.la/v1/users/%s/graphs/%s%s", param.Username, param.GraphID, paramsStr)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("could not get svg response : %s, %v", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		log.Printf("something went wrong to get svg response : %d %s", resp.StatusCode, url)
		w.WriteHeader(resp.StatusCode)
		return
	}

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("could not get svg response : %s, %v", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	extractData := string(byteArray)

	// flush to file
	tmpDirname := fmt.Sprintf("/tmp/pixela-svg/%s", time.Now().Format("2006-01-02"))
	tmpSvgFilePath := fmt.Sprintf("%s/%s.svg", tmpDirname, param.Hash)
	err = flushFile(tmpDirname, tmpSvgFilePath, extractData)
	if err != nil {
		log.Printf("failed to flush file : %s, %v", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// convert
	var size string
	if param.Mode == "short" {
		size = fmt.Sprintf("%sx%s", "220", "135")
	} else {
		size = fmt.Sprintf("%sx%s", "720", "135")
	}
	tmpPngDirname := fmt.Sprintf("/tmp/pixela-png/%s", time.Now().Format("2006-01-02"))
	tmpPngFilePath := fmt.Sprintf("%s/%s.png", tmpPngDirname, param.Hash)

	// make destination dir
	if _, err := os.Stat(tmpPngDirname); err != nil {
		if err := os.MkdirAll(tmpPngDirname, 0777); err != nil {
			log.Printf("failed to create dest dir : %s, %v", tmpPngDirname, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = exec.Command("convert", "-geometry", size, tmpSvgFilePath, tmpPngFilePath).Run()
	if err != nil {
		log.Printf("failed to run convert command : %s, %v", tmpSvgFilePath, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// response
	http.ServeFile(w, r, tmpPngFilePath)
}
