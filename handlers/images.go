package handlers

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"time"

	"cloud.google.com/go/storage"
	"github.com/Songmu/retry"
	"github.com/go-chi/chi"
)

type Target struct {
	ctx                context.Context
	originalRequestURI string
	githubID           string
	svgData            string
	tmpSvgFilePath     string
	goofysSvgFilePath  string
	tmpPngFilePath     string
	goofysPngFilePath  string
	size               string
	rotate             string
	transparent        bool
	date               time.Time
	pastdate           bool
}

func HandleImages(w http.ResponseWriter, r *http.Request) {
	t := &Target{
		githubID:           chi.URLParam(r, "githubID"),
		originalRequestURI: r.RequestURI,
		ctx:                r.Context(),
	}
	err := t.parseParams()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if t.pastdate {
		err := t.getPastSvgData()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		err := t.extractSvg()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	err = t.generatePng()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	http.ServeFile(w, r, t.tmpPngFilePath)
}

func (t *Target) parseParams() error {
	u, err := url.Parse(t.originalRequestURI)
	if err != nil {
		log.Printf("could not parse request URI : %v", err)
		return err
	}
	query := u.Query()

	width := "870"
	height := "155"
	t.rotate = "0"

	if query["width"] != nil {
		if regexp.MustCompile(`[^0-9]`).Match([]byte(query["width"][0])) {
			// invalid param, replace valid value
			width = "870"
		} else {
			width = query["width"][0]
		}
	}
	if query["height"] != nil {
		if regexp.MustCompile(`[^0-9]`).Match([]byte(query["height"][0])) {
			// invalid param, replace valid value
			height = "155"
		} else {
			height = query["height"][0]
		}
	}
	if query["rotate"] != nil {
		if regexp.MustCompile(`[^0-9]`).Match([]byte(query["rotate"][0])) {
			// invalid param, replace valid value
			t.rotate = "0"
		} else {
			t.rotate = query["rotate"][0]
		}
	}
	if query["date"] != nil {
		t.pastdate = true
		dateLayout := "20060102"
		t.date, err = time.Parse(dateLayout, query["date"][0])
		if err != nil {
			// invalid date param
			t.pastdate = false
			t.date = time.Now().Add(time.Duration(-10) * time.Minute)
		}
	} else {
		t.pastdate = false
		t.date = time.Now().Add(time.Duration(-10) * time.Minute)
	}
	t.size = fmt.Sprintf("%sx%s", width, height)

	if query["background"] != nil && query["background"][0] == "none" {
		t.transparent = true
	} else {
		t.transparent = false
	}
	return nil
}

func (t *Target) extractSvg() error {
	tmpDirname := fmt.Sprintf("/tmp/gg_svg/%s", t.date.Format("2006-01-02"))
	t.tmpSvgFilePath = fmt.Sprintf("%s/%s.svg", tmpDirname, t.githubID)
	goofysDirname := fmt.Sprintf("/var/www/grass-graph/goofys-s3/gg_svg/%s", t.date.Format("2006-01-02"))
	t.goofysSvgFilePath = fmt.Sprintf("%s/%s.svg", goofysDirname, t.githubID)

	if _, err := os.Stat(t.tmpSvgFilePath); err == nil {
		return nil
	}

	if _, err := os.Stat(t.goofysSvgFilePath); err == nil {
		err := copyFile(tmpDirname, t.tmpSvgFilePath, t.goofysSvgFilePath)
		if err != nil {
			return err
		}
		return nil
	}

	var byteArray []byte
	err := retry.Retry(5, 1*time.Second, func() error {
		url := fmt.Sprintf("https://github.com/%s", t.githubID)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("could not get github profile page response : %s, %v", url, err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 404 {
			log.Printf("could not get github profile page response : %s", url)
			return err
		}

		byteArray, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("could not get github profile page response : %s, %v", url, err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("could not get github profile page response (retry count exceeded): %v", err)
		return err
	}

	pageResponse := string(byteArray)

	// consider timeout response
	repexp := regexp.MustCompile(`too long to load`)
	if repexp.MatchString(pageResponse) {
		log.Printf("Probably, a timeout occurred. GitHub ID is %s .", t.githubID)
		return err
	}

	repexp = regexp.MustCompile(`^[\s\S]+<svg.+class="js-calendar-graph-svg">`)
	repcnd := `<svg xmlns="http://www.w3.org/2000/svg" width="870" height="155" class="js-calendar-graph-svg"><rect x="0" y="0" width="870" height="155" fill="white" stroke="none"/>`
	extractData := repexp.ReplaceAllString(pageResponse, repcnd)

	// Legend
	repexp = regexp.MustCompile(`dy="81" style="display: none;">Sat<\/text>[\s\S]+<\/g>[\s\S]+<\/svg>[.\s\S]+\z`)
	repcnd = `dy="81" style="display: none;">Sat</text><text x="675" y="125">Less</text><g transform="translate(709 , 15)"><rect class="day" width="11" height="11" x="0" y="99" fill="#eeeeee"/></g><g transform="translate(724 , 15)"><rect class="day" width="11" height="11" y="99" fill="#d6e685"/></g><g transform="translate(739 , 15)"><rect class="day" width="11" height="11" y="99" fill="#8cc665"/></g><g transform="translate(754 , 15)"><rect class="day" width="11" height="11" y="99" fill="#44a340"/></g><g transform="translate(769 , 15)"><rect class="day" width="11" height="11" y="99" fill="#1e6823"/></g><text x="788" y="125">More</text></g></svg>`
	extractData = repexp.ReplaceAllString(extractData, repcnd)

	// font-family
	repexp = regexp.MustCompile(`<text`)
	repcnd = `<text font-family="Helvetica"`
	extractData = repexp.ReplaceAllString(extractData, repcnd)
	t.svgData = extractData

	// output to file
	doneFlushFile := make(chan struct{}, 0)
	go func() error {
		err := flushFile(tmpDirname, t.tmpSvgFilePath, extractData)
		if err != nil {
			return err
		}
		err = flushFile(goofysDirname, t.goofysSvgFilePath, extractData)
		if err != nil {
			return err
		}
		close(doneFlushFile)
		return nil
	}()

	doneUpload := make(chan struct{}, 0)
	go func() error {
		err := t.uploadGcs(doneUpload)
		if err != nil {
			return err
		}
		return nil
	}()

	doneNotify := make(chan struct{}, 0)
	go func() {
		n := Notification{text: fmt.Sprintf("GitHub ID : %s's Grass-Graph Generated!!\nhttps://github.com/%s", t.githubID, t.githubID)}
		n.Notify()
		close(doneNotify)
	}()

	<-doneFlushFile
	<-doneUpload
	<-doneNotify

	return nil
}

func (t *Target) getPastSvgData() error {
	// get past graph data
	tmpDirname := fmt.Sprintf("/var/www/grass-graph/goofys-s3/gg_svg/%s", t.date.Format("2006-01-02"))
	t.tmpSvgFilePath = fmt.Sprintf("%s/%s.svg", tmpDirname, t.githubID)

	if _, err := os.Stat(t.tmpSvgFilePath); err == nil {
		// svg data file already exist
		return err
	}

	// download svg from gcs
	bucketname := "gg-on-a-know-home"
	objname := generateObjname(t)

	client, err := storage.NewClient(t.ctx)
	if err != nil {
		log.Printf("could not create gcs client : %v", err)
		return err
	}

	// GCS reader
	rc, err := client.Bucket(bucketname).Object(objname).NewReader(t.ctx)
	if err != nil {
		log.Printf("could not create gcs reader : %v", err)
		return err
	}
	defer rc.Close()

	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Printf("could not read gcs object : %v", err)
		return err
	}

	// output to file
	flushFile(tmpDirname, t.tmpSvgFilePath, ((string)(slurp)))

	return nil
}

func (t *Target) uploadGcs(uploadGcs chan struct{}) error {
	bucketname := "gg-on-a-know-home"
	objname := generateObjname(t)

	client, err := storage.NewClient(t.ctx)
	if err != nil {
		log.Printf("could not create gcs client : %v", err)
		return err
	}

	// GCS writer
	writer := client.Bucket(bucketname).Object(objname).NewWriter(t.ctx)
	writer.ContentType = "image/svg+xml"

	// upload : write object body
	if _, err := writer.Write(([]byte)(t.svgData)); err != nil {
		log.Printf("could not write object body : %v", err)
		return err
	}

	if err := writer.Close(); err != nil {
		log.Printf("could not close gcs writer : %v", err)
		return err
	}
	close(uploadGcs)
	return nil
}

func (t *Target) generatePng() error {

	tmpPngDirname := fmt.Sprintf("/tmp/gg_png/%s", t.date.Format("2006-01-02"))
	t.tmpPngFilePath = fmt.Sprintf("%s/%s-%s-%s-%s-%s.png", tmpPngDirname, t.githubID, t.size, t.rotate, t.transparent)
	goofysPngDirname := fmt.Sprintf("/var/www/grass-graph/goofys-s3/gg_png/%s", t.date.Format("2006-01-02"))
	t.goofysPngFilePath = fmt.Sprintf("%s/%s-%s-%s-%s-%s.png", goofysPngDirname, t.githubID, t.githubID, t.size, t.rotate, t.transparent)

	if _, err := os.Stat(t.tmpPngFilePath); err == nil {
		return nil
	}

	if _, err := os.Stat(t.goofysPngFilePath); err == nil {
		err := copyFile(tmpPngDirname, t.tmpPngFilePath, t.goofysPngFilePath)
		if err != nil {
			return err
		}
		return nil
	}

	// make destination dir
	if _, err := os.Stat(tmpPngDirname); err != nil {
		if err := os.MkdirAll(tmpPngDirname, 0777); err != nil {
			log.Printf("could not create direcotry : %v", err)
			return err
		}
	}

	if t.transparent {
		err := exec.Command("convert", "-geometry", t.size, "-rotate", t.rotate, "-transparent", "white", t.tmpSvgFilePath, t.tmpPngFilePath).Run()
		if err != nil {
			log.Printf("failed to run convert command : %s, %v", t.tmpSvgFilePath, err)
			return err
		}
	} else {
		err := exec.Command("convert", "-geometry", t.size, "-rotate", t.rotate, t.tmpSvgFilePath, t.tmpPngFilePath).Run()
		if err != nil {
			log.Printf("failed to run convert command : %s, %v", t.tmpSvgFilePath, err)
			return err
		}
	}
	err := copyFile(goofysPngDirname, t.goofysPngFilePath, t.tmpPngFilePath)
	if err != nil {
		return err
	}
	return nil
}

func generateObjname(t *Target) string {
	objpath := fmt.Sprintf("gg-svg-data/%s/%s/%s/%s", t.date.Format("2006"), t.date.Format("01"), t.date.Format("02"), t.githubID[0:1])
	return fmt.Sprintf("%s/%s_%s_graph.svg", objpath, t.githubID, t.date.Format("2006-01-02"))
}

func flushFile(dirname string, filepath string, data string) error {
	// make destination dir
	if _, err := os.Stat(dirname); err != nil {
		if err := os.MkdirAll(dirname, 0777); err != nil {
			log.Printf("could not crate tmp directory : %v", err)
			return err
		}
	}

	// output to file
	file, err := os.Create(filepath)
	if err != nil {
		log.Printf("could not crate svg file : %v", err)
		return err
	}
	defer file.Close()
	file.Write(([]byte)(data))
	return nil
}

func copyFile(dirname string, destFilepath string, sourceFilepath string) error {
	// make destination dir
	if _, err := os.Stat(dirname); err != nil {
		if err := os.MkdirAll(dirname, 0777); err != nil {
			log.Printf("could not crate tmp directory : %v", err)
			return err
		}
	}

	src, err := os.Open(sourceFilepath)
	if err != nil {
		log.Printf("could not open source svg file : %v", err)
		return err
	}
	defer src.Close()

	dst, err := os.Create(destFilepath)
	if err != nil {
		log.Printf("could not create destination svg file : %v", err)
		return err
	}
	defer dst.Close()

	// copy file
	_, err = io.Copy(dst, src)
	if err != nil {
		log.Printf("could not copy svg file : %v", err)
		return err
	}
	return nil
}
