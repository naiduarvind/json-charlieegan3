package collect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInstagram(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var content []byte
		var err error
		if strings.Contains(r.URL.Path, "/p/") {
			content, err = ioutil.ReadFile("instagram_response_post.html")
		} else {
			content, err = ioutil.ReadFile("instagram_response_profile.html")
		}
		if err != nil {
			t.Error(err)
		}
		fmt.Fprint(w, string(content))
	}))

	result, err := Instagram(localServer.URL, "charlieegan3")
	if err != nil {
		t.Error(err)
	}

	if result.Location != "Barbican Estate" {
		t.Error(result)
	}

	if strings.Contains(result.URL, "/p/BmCO0mAgC2h") == false {
		t.Error(result)
	}

	if result.Type != "photo" {
		t.Error(result)
	}
}
