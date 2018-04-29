package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
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
