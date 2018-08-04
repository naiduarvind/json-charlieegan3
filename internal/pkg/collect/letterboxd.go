package collect

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/pkg/errors"
)

type rssDocument struct {
	XMLName xml.Name `xml:"rss,omitempty" json:"rss,omitempty"`
	Channel *struct {
		XMLName xml.Name `xml:"channel,omitempty" json:"channel,omitempty"`
		Item    []*struct {
			XMLName xml.Name `xml:"item,omitempty" json:"item,omitempty"`
			Link    *struct {
				XMLName xml.Name `xml:"link,omitempty" json:"link,omitempty"`
				String  string   `xml:",chardata" json:",omitempty"`
			} `xml:"link,omitempty" json:"link,omitempty"`
			PubDate *struct {
				XMLName xml.Name `xml:"pubDate,omitempty" json:"pubDate,omitempty"`
				String  string   `xml:",chardata" json:",omitempty"`
			} `xml:"pubDate,omitempty" json:"pubDate,omitempty"`
			Title *struct {
				XMLName xml.Name `xml:"title,omitempty" json:"title,omitempty"`
				String  string   `xml:",chardata" json:",omitempty"`
			} `xml:"title,omitempty" json:"title,omitempty"`
		} `xml:"item,omitempty" json:"item,omitempty"`
	} `xml:"channel,omitempty" json:"channel,omitempty"`
}

// LatestFilm contains the wanted information for the latest film
type LatestFilm struct {
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Link      string    `json:"link"`
	Rating    string    `json:"rating"`
	Year      string    `json:"year"`
}

// Letterboxd returns the latest film in user's activity
func Letterboxd(host string, username string) (LatestFilm, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/rss", host, username))
	if err != nil {
		return LatestFilm{}, errors.Wrap(err, "get activities failed")
	}

	defer resp.Body.Close()

	var rss rssDocument
	err = xml.NewDecoder(resp.Body).Decode(&rss)
	if err != nil {
		return LatestFilm{}, errors.Wrap(err, "body unmarshal failed")
	}

	r := regexp.MustCompile(`(.*), (\d{4}) - (\S*)`)
	matches := r.FindSubmatch([]byte(rss.Channel.Item[0].Title.String))

	title := string(matches[1])
	year := string(matches[2])
	rating := string(matches[3])

	createdAt, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", rss.Channel.Item[0].PubDate.String)
	if err != nil {
		return LatestFilm{}, errors.Wrap(err, "failed to parse item date")
	}

	return LatestFilm{
		Title:     title,
		Year:      year,
		Rating:    rating,
		CreatedAt: createdAt,
		Link:      rss.Channel.Item[0].Link.String,
	}, nil
}
