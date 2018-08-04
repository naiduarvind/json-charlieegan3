package collect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLetterboxd(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("letterboxd_response.rss")
		if err != nil {
			t.Error(err)
		}
		fmt.Fprint(w, string(content))
	}))

	result, err := Letterboxd(localServer.URL, "charlieegan3")
	if err != nil {
		t.Error(err)
	}

	if result.Title != "Ready Player One" {
		t.Error(result)
	}
	if result.Year != "2018" {
		t.Error(result)
	}
	if strings.Contains(fmt.Sprintf("%v", result.CreatedAt), "2018-07-13 11:08:29") == false {
		t.Errorf("%v", result.CreatedAt)
	}
	if result.Rating != "★★★½" {
		t.Error(result.Rating)
	}
	if result.Link != "https://letterboxd.com/charlieegan3/film/ready-player-one/" {
		t.Error(result.Link)
	}
}
