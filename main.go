package main

import (
	"log"
	"net/http"
	"time"

	"github.com/a-know/grass-graph/handlers"
	"github.com/fukata/golang-stats-api-handler"
	"github.com/go-chi/chi"
)

const location = "Asia/Tokyo"

func main() {
	r := chi.NewRouter()

	// timezone
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc

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

	// for monitoring
	r.Get("/api/stats", stats_api.Handler)

	r.Post("/knock", handlers.HandleKnock)

	// for Pixela SVG convert to PNG
	// /pixela/convert?username=a-know&graphID=test-graph&date=yyyyMMdd&mode=short&stage=dev&hash=xxxxx
	r.Get("/pixela/convert", handlers.HandleSVGConvert)

	log.Printf("grass-graph started.")
	http.ListenAndServe(":8080", r)
}
