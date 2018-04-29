package handlers

import (
	"io/ioutil"
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
			//TODO logger
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			//TODO logger
		}

		var ns = template.New("template")
		t.templ, _ = ns.Parse(string(data))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	t.templ.Execute(w, data)
}
