package collect

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type event struct {
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	Payload   struct {
		Commits []struct {
			Message string `json:"message"`
			URL     string `json:"url"`
		} `json:"commits"`
	}
	Repo struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"repo"`
}

// LatestCommit stores the message, time and repo of the user's latest commit
type LatestCommit struct {
	Message string `json:"message"`
	URL     string `json:"url"`
	Repo    struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"repo"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedAtString string    `json:"created_at_string"`
}

// Collect returns a user's latest commit and project
func (l *LatestCommit) Collect(host string, username string) error {
	resp, err := http.Get(fmt.Sprintf("%s/users/%s/events", host, username))
	if err != nil {
		return errors.Wrap(err, "GitHub get failed")
	}

	defer resp.Body.Close()

	var events []event
	err = json.NewDecoder(resp.Body).Decode(&events)
	if err != nil {
		return errors.Wrap(err, "GitHub body unmarshal failed")
	}

	var pushes []event
	for _, event := range events {
		if event.Type == "PushEvent" {
			pushes = append(pushes, event)
			break
		}
	}

	if len(pushes) < 1 {
		return errors.New("GitHub response contained no pushes")
	}

	latestPush := pushes[0]

	createdAt, err := time.Parse(time.RFC3339, latestPush.CreatedAt)
	if err != nil {
		return errors.Wrap(err, "GitHub latest event time parsing failed")
	}

	l.CreatedAt = createdAt
	l.Repo = latestPush.Repo
	commit := latestPush.Payload.Commits[len(latestPush.Payload.Commits)-1]
	l.URL = commit.URL
	l.Message = commit.Message

	return nil
}
