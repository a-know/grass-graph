package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/go-chi/chi"
)

type Target struct {
	originalRequestURI string
	githubID           string
	svgData            string
	tmpSvgFilePath     string
	tmpPngFilePath     string
	size               string
	rotate             string
	transparent        bool
}

func HandleImages(w http.ResponseWriter, r *http.Request) {
	t := &Target{githubID: chi.URLParam(r, "githubID"), originalRequestURI: r.RequestURI}
	t.parseParams()
	t.extractSvg()
	t.generatePng()

	http.ServeFile(w, r, t.tmpPngFilePath)

	// contributions_info := regexp.MustCompile("<span class=\"contrib-number\">(.+)</span>")
	// group := assined.FindSubmatch(byteArray)
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
	t.size = fmt.Sprintf("%sx%s", width, height)

	if query["background"] != nil && query["background"][0] == "none" {
		t.transparent = true
	} else {
		t.transparent = false
	}
}

func (t *Target) extractSvg() {
	tmpDirname := fmt.Sprintf("tmp/gg_svg/%s", time.Now().Format("2006-01-02"))
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
}

func (t *Target) generatePng() {

	tmpPngDirname := fmt.Sprintf("tmp/gg_png/%s", time.Now().Format("2006-01-02"))
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
