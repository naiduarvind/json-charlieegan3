package collect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStrava(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("strava_response.json")
		if err != nil {
			t.Error(err)
		}
		fmt.Fprint(w, string(content))
	}))

	result, err := Strava(localServer.URL)

	if err != nil {
		t.Error(err)
	}

	if result.Name != "Evaporate" {
		t.Error(result)
	}
	if result.Distance != 4231.5 {
		t.Error(result)
	}
	if result.MovingTime != 1470 {
		t.Error(result)
	}
	if result.AverageHeartrate != 142.9 {
		t.Error(result)
	}
	if result.Type != "Run" {
		t.Error(result)
	}
	if result.URL != "https://www.strava.com/activities/1748439744" {
		t.Error(result)
	}
}
