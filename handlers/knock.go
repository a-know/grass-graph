package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Notification struct {
	text string
}

func HandleKnock(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ua := r.Form.Get("user_agent")
	lang := r.Form.Get("language")
	admin := r.Form.Get("admin")

	if admin != "true" {
		n := &Notification{text: fmt.Sprintf("Visitor Incoming!!\nUA : %s\nLanguage : %s", ua, lang)}
		err := n.Notify()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (n *Notification) Notify() error {
	name := os.Getenv("SLACK_BOT_NAME")
	channel := os.Getenv("SLACK_CHANNEL_NAME")

	jsonStr := `{"channel":"` + channel + `","username":"` + name + `","text":"` + n.text + `"}`

	req, err := http.NewRequest(
		"POST",
		os.Getenv("SLACK_WEBHOOK_URL"),
		bytes.NewBuffer([]byte(jsonStr)),
	)

	if err != nil {
		log.Printf("could not request to slack webhook : %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("could not request to slack webhook : %v", err)
		return err
	}

	defer resp.Body.Close()

	return nil
}
