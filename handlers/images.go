package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/go-chi/chi"
)

type Target struct {
	githubID string
	svgData  string
}

func HandleImages(w http.ResponseWriter, r *http.Request) {
	t := &Target{githubID: chi.URLParam(r, "githubID")}
	t.extractSvg()

	// contributions_info := regexp.MustCompile("<span class=\"contrib-number\">(.+)</span>")
	// group := assined.FindSubmatch(byteArray)
}

func (t *Target) extractSvg() {
	tmpDirname := fmt.Sprintf("tmp/gg_svg/%s", time.Now().Format("2006-01-02"))
	tmpFilename := fmt.Sprintf("%s/%s.svg", tmpDirname, t.githubID)

	if _, err := os.Stat(tmpFilename); err == nil {
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
	repcnd = `dy="81" style="display: none;">Sat</text><text x="535" y="110">Less</text><g transform="translate(569 , 0)"><rect class="day" width="11" height="11" y="99" fill="#eeeeee"/></g><g transform="translate(584 , 0)"><rect class="day" width="11" height="11" y="99" fill="#d6e685"/></g><g transform="translate(599 , 0)"><rect class="day" width="11" height="11" y="99" fill="#8cc665"/></g><g transform="translate(614 , 0)"><rect class="day" width="11" height="11" y="99" fill="#44a340"/></g><g transform="translate(629 , 0)"><rect class="day" width="11" height="11" y="99" fill="#1e6823"/></g><text x="648" y="110">More</text></g></svg>`
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
	file, err := os.Create(tmpFilename)
	if err != nil {
		// TODO logger
	}
	defer file.Close()
	file.Write(([]byte)(extractData))
}
