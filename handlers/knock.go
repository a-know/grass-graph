package handlers

import (
	"bytes"
	"fmt"
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
		n.Notify()
	}
}

func (n *Notification) Notify() {
	name := os.Getenv("SLACK_BOT_NAME")
	channel := os.Getenv("SLACK_CHANNEL_NAME")

	jsonStr := `{"channel":"` + channel + `","username":"` + name + `","text":"` + n.text + `"}`

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
