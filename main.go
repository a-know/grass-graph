package main

import (
	"fmt"
	"net/http"
	"os"
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

	r.Get("/images/{githubID}.png", handlers.HandleImages)

	// for monitoring
	r.Get("/heartbeat", handlers.HandleHeartbeat)
	r.Get("/api/stats", stats_api.Handler)

	// r.Post("/knock", handlers.HandleKnock)

	// for Pixela SVG convert to PNG
	// /pixela/convert?username=a-know&graphID=test-graph&date=yyyyMMdd&mode=short&stage=dev&hash=xxxxx
	r.Get("/pixela/convert", handlers.HandleSVGConvert)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
