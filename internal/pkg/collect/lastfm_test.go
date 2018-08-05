package collect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLastFm(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("lastfm_response.json")
		if err != nil {
			t.Error(err)
		}
		fmt.Fprint(w, string(content))
	}))

	var latestTrack LatestTrack
	err := latestTrack.Collect(localServer.URL, "charlieegan3", "KEY")
	if err != nil {
		t.Error(err)
	}

	if latestTrack.Name != "The Trip" {
		t.Error(latestTrack.Name)
	}
	if latestTrack.Link != "https://www.last.fm/music/Still+Corners/_/The+Trip" {
		t.Error(latestTrack.Link)
	}
	if strings.Contains(latestTrack.ProfileLink, "/user/charlieegan3") == false {
		t.Error(latestTrack.ProfileLink)
	}
	if latestTrack.Artist != "Still Corners" {
		t.Error(latestTrack.Artist)
	}
	if latestTrack.CreatedAt.String() != "2018-08-04 22:11:51 +0100 BST" {
		t.Error(latestTrack.CreatedAt)
	}
}
