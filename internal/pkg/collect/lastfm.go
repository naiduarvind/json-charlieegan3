package collect

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type response struct {
	Recenttracks struct {
		Track []struct {
			Artist struct {
				Text string `json:"#text"`
			} `json:"artist"`
			Date struct {
				Uts string `json:"uts"`
			} `json:"date"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"track"`
	} `json:"recenttracks"`
}

// LatestTrack wraps the data for the most recent lastfm track
type LatestTrack struct {
	Link            string    `json:"link"`
	ProfileLink     string    `json:"profile"`
	Name            string    `json:"name"`
	Artist          string    `json:"artist"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedAtString string    `json:"created_at_string"`
}

// Collect gets the latest last fm track for a given user
// baseURL https://ws.audioscrobbler.com
func (l *LatestTrack) Collect(baseURL string, username string, apiKey string) error {
	resp, err := http.Get(fmt.Sprintf("%s/2.0/?method=user.getrecenttracks&user=%s&api_key=%s&format=json", baseURL, username, apiKey))
	if err != nil {
		return errors.Wrap(err, "get recent tracks failed")
	}
	defer resp.Body.Close()

	var data response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return errors.Wrap(err, "body unmarshal failed")
	}

	track := data.Recenttracks.Track[0]

	if track.Date.Uts == "" {
		track.Date.Uts = fmt.Sprintf("%d", time.Now().Unix())
	}
	uts, err := strconv.ParseInt(track.Date.Uts, 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse track unix timestamp")
	}

	l.Link = track.URL
	l.ProfileLink = fmt.Sprintf("%s/user/%s", baseURL, username)
	l.Name = track.Name
	l.Artist = track.Artist.Text
	l.CreatedAt = time.Unix(uts, 0)

	return nil
}
