package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/a-know/grass-graph-go/handlers"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	t := &handlers.TemplateHandler{Filename: "index.html", Assets: Assets}
	r.Get("/", t.HandleTemplate)

	css := &handlers.AssetsHandler{Kind: "css"}
	r.Get("/css/*", css.HandleAssets)

	js := &handlers.AssetsHandler{Kind: "js"}
	r.Get("/js/*", js.HandleAssets)

	fonts := &handlers.AssetsHandler{Kind: "fonts"}
	r.Get("/fonts/*", fonts.HandleAssets)

	images := &handlers.AssetsHandler{Kind: "images"}
	r.Get("/images/*", images.HandleAssets)

	plugins := &handlers.AssetsHandler{Kind: "plugins"}
	r.Get("/plugins/*", plugins.HandleAssets)

	r.Post("/knock", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		ua := r.Form.Get("user_agent")
		lang := r.Form.Get("language")
		admin := r.Form.Get("admin")

		if admin != "true" {
			name := os.Getenv("SLACK_BOT_NAME")
			text := fmt.Sprintf("Visitor Incoming!!\nUA : %s\nLanguage : %s", ua, lang)
			channel := os.Getenv("SLACK_CHANNEL_NAME")

			jsonStr := `{"channel":"` + channel + `","username":"` + name + `","text":"` + text + `"}`

			req, err := http.NewRequest(
				"POST",
				os.Getenv("SLACK_WEBHOOK_URL"),
				bytes.NewBuffer([]byte(jsonStr)),
			)

			if err != nil {
				fmt.Print(err)
			}

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Print(err)
			}

			fmt.Print(resp)
			defer resp.Body.Close()
		}
	})
	http.ListenAndServe(":8080", r)
}
