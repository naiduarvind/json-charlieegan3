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
	CreatedAt        time.Time `json:"created_at"`
	Type             string    `json:"type"`
}

// Strava returns details about the most recent strava activity
func Strava(host string) (LatestActivity, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v3/athlete/activities", host))
	if err != nil {
		return LatestActivity{}, errors.Wrap(err, "get activities failed")
	}

	defer resp.Body.Close()

	var activities []activity
	err = json.NewDecoder(resp.Body).Decode(&activities)
	if err != nil {
		return LatestActivity{}, errors.Wrap(err, "body unmarshal failed")
	}

	if len(activities) == 0 {
		return LatestActivity{}, errors.New("no activities found")
	}

	activity := activities[0]
	createdAt, err := time.Parse(time.RFC3339, activity.StartDate)
	if err != nil {
		return LatestActivity{}, errors.Wrap(err, "latest activity time parsing failed")
	}

	return LatestActivity{
		AverageHeartrate: activity.AverageHeartrate,
		Distance:         activity.Distance,
		MovingTime:       activity.MovingTime,
		Name:             activity.Name,
		Type:             activity.Type,
		CreatedAt:        createdAt,
		URL:              fmt.Sprintf("https://www.strava.com/activities/%d", activity.ID),
	}, nil
}
