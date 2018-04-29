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

	r.Get("/knock", func(w http.ResponseWriter, r *http.Request) {
		name := os.Getenv("SLACK_BOT_NAME")
		text := "Visitor Incoming!!\nUA : hoge\nLanguage : fuga"
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
	})
	http.ListenAndServe(":8080", r)
}
