package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"text/template"

	"github.com/jessevdk/go-assets"
)

type TemplateHandler struct {
	once     sync.Once
	Filename string
	templ    *template.Template
	Assets   *assets.FileSystem
}

// Handling HTTP Request
func (t *TemplateHandler) HandleTemplate(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		f, err := t.Assets.Open(filepath.Join("/public", t.Filename))
		if err != nil {
			log.Printf("could not open assets file : %v", err)
			w.WriteHeader(http.StatusNotFound)
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			log.Printf("could not read assets file : %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		var ns = template.New("template")
		t.templ, _ = ns.Parse(string(data))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	t.templ.Execute(w, data)
}
