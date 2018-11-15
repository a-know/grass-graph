package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// ConvertParam :
type ConvertParam struct {
	Username string `json:"username"`
	GraphID  string `json:"graphID"`
	Date     string `json:"date"`
	Mode     string `json:"mode"`
	Stage    string `json:"stage"`
	Hash string `json:"hash"`
}

// HandleSVGConvert :
func HandleSVGConvert(w http.ResponseWriter, r *http.Request) {
	// unmarshal
	var param ConvertParam
	err = json.Unmarshal(b, &param)
	if err != nil {
		log.Printf("Failed to unmarshal params: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// treat params
	var paramsStr string
	if param.Date != "" {
		paramsStr = fmt.Sprintf("?date=%s", param.Date)
	}
	if param.Mode != "" {
		paramsStr = fmt.Sprintf("%s&mode=%s", paramsStr, param.Mode)
	}

	// get svg
	var url string
	if param.Stage == "dev" {
		url = fmt.Sprintf("https://www-dev-215102.appspot.com/v1/users/%s/graphs/%s%s", param.Username, param.GraphID, paramsStr)
	} else {
		url = fmt.Sprintf("https://pixela/v1/users/%s/graphs/%s%s", param.Username, param.GraphID, paramsStr)
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

	byteArray, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("could not get svg response : %s, %v", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	extractData := string(byteArray)

	// flush to file
	tmpDirname := fmt.Sprintf("/tmp/pixela-svg/%s", t.date.Format("2006-01-02"))
	tmpSvgFilePath = fmt.Sprintf("%s/%s.svg", tmpDirname, param.Hash)
	err := flushFile(tmpDirname, t.tmpSvgFilePath, extractData)
	if err != nil {
		return err
	}

	// convert
	var size string
	if param.Mode == "short" {
		size = fmt.Sprintf("%sx%s", "220", "135")
	} else {
		size = fmt.Sprintf("%sx%s", "720", "135")
	}
	tmpPngDirname := fmt.Sprintf("/tmp/pixela-png/%s", t.date.Format("2006-01-02"))
	tmpPngFilePath = fmt.Sprintf("%s/%s.png", tmpPngDirname, param.Hash)
	err := exec.Command("convert", "-geometry", t.size, tmpSvgFilePath, tmpPngFilePath).Run()
	if err != nil {
		log.Printf("failed to run convert command : %s, %v", t.tmpSvgFilePath, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// response
	http.ServeFile(w, r, tmpPngFilePath)
}
