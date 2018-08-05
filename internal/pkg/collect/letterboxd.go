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

// Collect returns the latest film in user's activity
func (l *LatestFilm) Collect(host string, username string) error {
	resp, err := http.Get(fmt.Sprintf("%s/%s/rss", host, username))
	if err != nil {
		return errors.Wrap(err, "get activities failed")
	}

	defer resp.Body.Close()

	var rss rssDocument
	err = xml.NewDecoder(resp.Body).Decode(&rss)
	if err != nil {
		return errors.Wrap(err, "body unmarshal failed")
	}

	r := regexp.MustCompile(`(.*), (\d{4}) - (\S*)`)
	matches := r.FindSubmatch([]byte(rss.Channel.Item[0].Title.String))

	createdAt, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", rss.Channel.Item[0].PubDate.String)
	if err != nil {
		return errors.Wrap(err, "failed to parse item date")
	}

	l.Title = string(matches[1])
	l.Year = string(matches[2])
	l.Rating = string(matches[3])
	l.CreatedAt = createdAt
	l.Link = rss.Channel.Item[0].Link.String

	return nil
}
