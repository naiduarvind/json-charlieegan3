package collect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTwitter(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("twitter_response.json")
		if err != nil {
			t.Error(err)
		}
		fmt.Fprint(w, string(content))
	}))

	result, err := Twitter(localServer.URL, []string{"t", "t", "t", "t"})
	if err != nil {
		t.Error(err)
	}

	if result.Text != "just another test" {
		t.Error(result.Text)
	}
	if result.Link != "https://twitter.com/oauth_dancer/status/240558470661799936" {
		t.Error(result.Link)
	}
	if fmt.Sprintf("%v", result.CreatedAt) != "2012-08-28 21:16:23 +0000 +0000" {
		t.Error(result.CreatedAt)
	}
	if result.Location != "Berkeley" {
		t.Error(result.Location)
	}
}
