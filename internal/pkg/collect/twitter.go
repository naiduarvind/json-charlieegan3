package collect

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

// LatestTweet wraps the required data for a tweet
type LatestTweet struct {
	Text      string    `json:"text"`
	Link      string    `json:"link"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

// Twitter returns the latest tweet for the requesting user
// Use https://api.twitter.com/1.1 as the baseURL
func Twitter(baseURL string, credentials []string) (LatestTweet, error) {
	api := anaconda.NewTwitterApiWithCredentials(credentials[0], credentials[1], credentials[2], credentials[3])
	api.SetBaseUrl(baseURL)

	params := url.Values{}
	params.Set("include_entities", "false")
	data, err := api.GetHomeTimeline(params)
	if err != nil {
		return LatestTweet{}, err
	}

	createdAt, err := data[0].CreatedAtTime()
	if err != nil {
		return LatestTweet{}, err
	}

	return LatestTweet{
		Text:      data[0].Text,
		CreatedAt: createdAt,
		Location:  data[0].Place.Name,
		Link:      fmt.Sprintf("https://twitter.com/%s/status/%s", data[0].User.ScreenName, data[0].IdStr),
	}, nil
}
