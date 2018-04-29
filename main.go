package main

import (
	"net/http"

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

	r.Get("/images/{githubID}.png", handlers.HandleImages)

	images := &handlers.AssetsHandler{Kind: "images"}
	r.Get("/images/*", images.HandleAssets)

	plugins := &handlers.AssetsHandler{Kind: "plugins"}
	r.Get("/plugins/*", plugins.HandleAssets)

	r.Post("/knock", handlers.HandleKnock)

	http.ListenAndServe(":8080", r)
}
