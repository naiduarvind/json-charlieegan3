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

	result, err := LastFm(localServer.URL, "charlieegan3", "KEY")
	if err != nil {
		t.Error(err)
	}

	if result.Name != "The Trip" {
		t.Error(result.Name)
	}
	if result.Link != "https://www.last.fm/music/Still+Corners/_/The+Trip" {
		t.Error(result.Link)
	}
	if strings.Contains(result.ProfileLink, "/user/charlieegan3") == false {
		t.Error(result.ProfileLink)
	}
	if result.Artist != "Still Corners" {
		t.Error(result.Artist)
	}
	if result.CreatedAt.String() != "2018-08-04 22:11:51 +0100 BST" {
		t.Error(result.CreatedAt)
	}
}
