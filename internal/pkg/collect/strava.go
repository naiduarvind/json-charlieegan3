package collect

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type activity struct {
	AverageHeartrate float64 `json:"average_heartrate"`
	ID               int64   `json:"id"`
	Distance         float64 `json:"distance"`
	MovingTime       int64   `json:"moving_time"`
	Name             string  `json:"name"`
	StartDate        string  `json:"start_date"`
	Type             string  `json:"type"`
}

// LatestActivity wraps deta about the latest activity
type LatestActivity struct {
	AverageHeartrate float64   `json:"average_heartrate"`
	URL              string    `json:"url"`
	Distance         float64   `json:"distance"`
	MovingTime       int64     `json:"moving_time"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedAtString  string    `json:"created_at_string"`
}

// Collect returns details about the most recent strava activity
// host https://www.strava.com
func (l *LatestActivity) Collect(host string, token string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v3/athlete/activities", host), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	if err != nil {
		return errors.Wrap(err, "build request failed")
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "get activities failed")
	}

	defer resp.Body.Close()

	var activities []activity
	err = json.NewDecoder(resp.Body).Decode(&activities)
	if err != nil {
		return errors.Wrap(err, "body unmarshal failed")
	}

	if len(activities) == 0 {
		return errors.New("no activities found")
	}

	activity := activities[0]
	createdAt, err := time.Parse(time.RFC3339, activity.StartDate)
	if err != nil {
		return errors.Wrap(err, "latest activity time parsing failed")
	}

	l.AverageHeartrate = activity.AverageHeartrate
	l.Distance = activity.Distance
	l.MovingTime = activity.MovingTime
	l.Name = activity.Name
	l.Type = activity.Type
	l.CreatedAt = createdAt
	l.URL = fmt.Sprintf("https://www.strava.com/activities/%d", activity.ID)

	return nil
}
