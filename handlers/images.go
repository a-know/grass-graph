package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"time"

	"cloud.google.com/go/storage"
	"github.com/go-chi/chi"
)

type Target struct {
	ctx                context.Context
	originalRequestURI string
	githubID           string
	svgData            string
	tmpSvgFilePath     string
	tmpPngFilePath     string
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
	t.parseParams()

	if t.pastdate {
		t.getPastSvgData()
	} else {
		t.extractSvg()
	}

	t.generatePng()

	http.ServeFile(w, r, t.tmpPngFilePath)
}

func (t *Target) parseParams() {
	u, err := url.Parse(t.originalRequestURI)
	if err != nil {
		// TODO logger
		panic(err)
	}
	query := u.Query()

	width := "720"
	height := "135"
	t.rotate = "0"
	// TODO validate
	if query["width"] != nil {
		width = query["width"][0]
	}
	if query["height"] != nil {
		height = query["height"][0]
	}
	if query["rotate"] != nil {
		t.rotate = query["rotate"][0]
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
}

func (t *Target) extractSvg() {
	tmpDirname := fmt.Sprintf("tmp/gg_svg/%s", t.date.Format("2006-01-02"))
	t.tmpSvgFilePath = fmt.Sprintf("%s/%s.svg", tmpDirname, t.githubID)

	if _, err := os.Stat(t.tmpSvgFilePath); err == nil {
		return
	}

	// TODO retry

	url := fmt.Sprintf("https://github.com/%s", t.githubID)
	resp, err := http.Get(url)
	if err != nil {
		// TODO logger
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		// TODO logger
		panic(err)
	}

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// TODO logger
		panic(err)
	}
	pageResponse := string(byteArray)

	repexp := regexp.MustCompile(`^[\s\S]+<svg.+class="js-calendar-graph-svg">`)
	repcnd := `<svg xmlns="http://www.w3.org/2000/svg" width="720" height="135" class="js-calendar-graph-svg"><rect x="0" y="0" width="720" height="135" fill="white" stroke="none"/>`
	extractData := repexp.ReplaceAllString(pageResponse, repcnd)

	// Legend
	repexp = regexp.MustCompile(`dy="81" style="display: none;">Sat<\/text>[\s\S]+<\/g>[\s\S]+<\/svg>[.\s\S]+\z`)
	repcnd = `dy="81" style="display: none;">Sat</text><text x="535" y="110">Less</text><g transform="translate(569 , 0)"><rect class="day" width="11" height="11" x="0" y="99" fill="#eeeeee"/></g><g transform="translate(584 , 0)"><rect class="day" width="11" height="11" y="99" fill="#d6e685"/></g><g transform="translate(599 , 0)"><rect class="day" width="11" height="11" y="99" fill="#8cc665"/></g><g transform="translate(614 , 0)"><rect class="day" width="11" height="11" y="99" fill="#44a340"/></g><g transform="translate(629 , 0)"><rect class="day" width="11" height="11" y="99" fill="#1e6823"/></g><text x="648" y="110">More</text></g></svg>`
	extractData = repexp.ReplaceAllString(extractData, repcnd)

	// font-family
	repexp = regexp.MustCompile(`<text`)
	repcnd = `<text font-family="Helvetica"`
	extractData = repexp.ReplaceAllString(extractData, repcnd)
	t.svgData = extractData

	// make destination dir
	if _, err := os.Stat(tmpDirname); err != nil {
		if err := os.MkdirAll(tmpDirname, 0777); err != nil {
			// TODO logger
		}
	}

	// output to file
	file, err := os.Create(t.tmpSvgFilePath)
	if err != nil {
		// TODO logger
	}
	defer file.Close()
	file.Write(([]byte)(extractData))

	t.uploadGcs()
}

func (t *Target) getPastSvgData() {
	// get past graph data
	tmpDirname := fmt.Sprintf("tmp/gg_svg/%s", t.date.Format("2006-01-02"))
	t.tmpSvgFilePath = fmt.Sprintf("%s/%s.svg", tmpDirname, t.githubID)

	if _, err := os.Stat(t.tmpSvgFilePath); err == nil {
		// svg data file already exist
		return
	}

	// download svg from gcs
	bucketname := "gg-on-a-know-home"
	objpath := fmt.Sprintf("gg-svg-data/%s/%s/%s/%s", t.date.Format("2006"), t.date.Format("01"), t.date.Format("02"), t.githubID[0:1])
	objname := fmt.Sprintf("%s/%s_%s_graph.svg", objpath, t.githubID, t.date.Format("2006-01-02"))

	client, err := storage.NewClient(t.ctx)
	if err != nil {
		// TODO logger
		panic(err)
	}

	// GCS reader
	rc, err := client.Bucket(bucketname).Object(objname).NewReader(t.ctx)
	if err != nil {
		// TODO logger
		panic(err)
	}
	defer rc.Close()

	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		// TODO logger
		panic(err)
	}

	// make destination dir
	if _, err := os.Stat(tmpDirname); err != nil {
		if err := os.MkdirAll(tmpDirname, 0777); err != nil {
			// TODO logger
		}
	}
	// output to file
	file, err := os.Create(t.tmpSvgFilePath)
	if err != nil {
		// TODO logger
	}
	defer file.Close()
	file.Write(([]byte)(slurp))
}

func (t *Target) uploadGcs() {
	bucketname := "gg-on-a-know-home"
	objpath := fmt.Sprintf("gg-svg-data/%s/%s/%s/%s", t.date.Format("2006"), t.date.Format("01"), t.date.Format("02"), t.githubID[0:1])
	objname := fmt.Sprintf("%s/%s_%s_graph.svg", objpath, t.githubID, t.date.Format("2006-01-02"))

	client, err := storage.NewClient(t.ctx)
	if err != nil {
		// TODO logger
		panic(err)
	}

	// GCS writer
	writer := client.Bucket(bucketname).Object(objname).NewWriter(t.ctx)
	writer.ContentType = "image/svg+xml"

	// upload : write object body
	if _, err := writer.Write(([]byte)(t.svgData)); err != nil {
		// TODO logger
		panic(err)
	}

	if err := writer.Close(); err != nil {
		// TODO logger
		panic(err)
	}
}

func (t *Target) generatePng() {

	tmpPngDirname := fmt.Sprintf("tmp/gg_png/%s", t.date.Format("2006-01-02"))
	t.tmpPngFilePath = fmt.Sprintf("%s/%s.png", tmpPngDirname, t.githubID)
	// make destination dir
	if _, err := os.Stat(tmpPngDirname); err != nil {
		if err := os.MkdirAll(tmpPngDirname, 0777); err != nil {
			// TODO logger
			panic(err)
		}
	}

	if t.transparent {
		err := exec.Command("convert", "-geometry", t.size, "-rotate", t.rotate, "-transparent", "white", t.tmpSvgFilePath, t.tmpPngFilePath).Run()
		if err != nil {
			// TODO logger
			panic(err)
		}
	} else {
		err := exec.Command("convert", "-geometry", t.size, "-rotate", t.rotate, t.tmpSvgFilePath, t.tmpPngFilePath).Run()
		if err != nil {
			// TODO logger
			panic(err)
		}
	}
}
