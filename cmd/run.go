package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	humanize "github.com/dustin/go-humanize"

	c "github.com/charlieegan3/json-charlieegan3/internal/pkg/collect"
	"github.com/pkg/errors"
)

type status struct {
	Tweet    c.LatestTweet    `json:"tweet"`
	Post     c.LatestPost     `json:"post"`
	Activity c.LatestActivity `json:"activity"`
	Film     c.LatestFilm     `json:"film"`
	Track    c.LatestTrack    `json:"track"`
	Commit   c.LatestCommit   `json:"commit"`
}

func (s *status) setCreatedAtStrings() {
	s.Tweet.CreatedAtString = compactHumanizeTime(s.Tweet.CreatedAt)
	s.Post.CreatedAtString = compactHumanizeTime(s.Post.CreatedAt)
	s.Activity.CreatedAtString = compactHumanizeTime(s.Activity.CreatedAt)
	s.Film.CreatedAtString = compactHumanizeTime(s.Film.CreatedAt)
	s.Track.CreatedAtString = compactHumanizeTime(s.Track.CreatedAt)
	s.Commit.CreatedAtString = compactHumanizeTime(s.Commit.CreatedAt)
}

func (s *status) fetchCurrent() error {
	url := fmt.Sprintf("%s/%s/%s", os.Getenv("AWS_BUCKET_HOST"), os.Getenv("AWS_BUCKET"), os.Getenv("STATUS_KEY"))

	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "current status get failed")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return errors.Wrap(err, "current status unmarshal failed")
	}

	return nil
}

func (s *status) fetchNew(previousStatus status) {
	username := os.Getenv("USERNAME")

	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		defer wg.Done()
		twitterCredentials := strings.Split(os.Getenv("TWITTER_CREDENTIALS"), ",")
		err := s.Tweet.Collect("https://api.twitter.com/1.1", twitterCredentials)
		if err != nil {
			fmt.Println(errors.Wrap(err, "twitter error"))
			s.Tweet = previousStatus.Tweet
		}
	}()

	go func() {
		defer wg.Done()
		err := s.Post.Collect("https://instagram.com", username)
		if err != nil {
			fmt.Println(errors.Wrap(err, "instagram error"))
			s.Post = previousStatus.Post
		}
	}()

	go func() {
		defer wg.Done()
		err := s.Activity.Collect("https://www.strava.com", os.Getenv("STRAVA_TOKEN"))
		if err != nil {
			fmt.Println(errors.Wrap(err, "strava error"))
			s.Activity = previousStatus.Activity
		}
	}()

	go func() {
		defer wg.Done()
		err := s.Film.Collect("https://letterboxd.com", username)
		if err != nil {
			fmt.Println(errors.Wrap(err, "letterboxd error"))
			s.Film = previousStatus.Film
		}
	}()

	go func() {
		defer wg.Done()
		err := s.Track.Collect("https://ws.audioscrobbler.com", username, os.Getenv("LASTFM_KEY"))
		if err != nil {
			fmt.Println(errors.Wrap(err, "lastfm error"))
			s.Track = previousStatus.Track
		}
	}()

	go func() {
		defer wg.Done()
		err := s.Commit.Collect("https://api.github.com", username)
		if err != nil {
			fmt.Println(errors.Wrap(err, "github error"))
			s.Commit = previousStatus.Commit
		}
	}()

	wg.Wait()
}

func upload(statusJSON string) error {
	conf := aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))}
	sess := session.New(&conf)

	statusKey := os.Getenv("STATUS_KEY")

	s3Service := s3manager.NewUploader(sess)

	_, err := s3Service.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(os.Getenv("STATUS_KEY")),
		Body:   strings.NewReader(statusJSON),
	})
	if err != nil {
		return errors.Wrap(err, "s3 upload error")
	}

	cloudfrontService := cloudfront.New(sess)
	distribution := os.Getenv("AWS_DISTRIBUTION")
	callerReference := "json-charlieegan3-go"
	path := fmt.Sprintf("/%s", statusKey)
	pathQuantity := int64(1)

	input := &cloudfront.CreateInvalidationInput{
		DistributionId: &distribution,
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: &callerReference,
			Paths: &cloudfront.Paths{
				Items: []*string{
					&path,
				},
				Quantity: &pathQuantity,
			},
		},
	}

	_, err = cloudfrontService.CreateInvalidation(input)
	if err != nil {
		return errors.Wrap(err, "cloudfront invalidation error")
	}

	return nil
}

func compactHumanizeTime(time time.Time) string {
	humanReadable := humanize.Time(time)

	humanReadable = strings.Replace(humanReadable, " year", "yr", 1)
	humanReadable = strings.Replace(humanReadable, " month", "mth", 1)
	humanReadable = strings.Replace(humanReadable, " weeks", "w", 1)
	humanReadable = strings.Replace(humanReadable, " week", "w", 1)
	humanReadable = strings.Replace(humanReadable, " days", "d", 1)
	humanReadable = strings.Replace(humanReadable, " day", "d", 1)
	humanReadable = strings.Replace(humanReadable, " hours", "h", 1)
	humanReadable = strings.Replace(humanReadable, " hour", "h", 1)
	humanReadable = strings.Replace(humanReadable, " minutes", "m", 1)
	humanReadable = strings.Replace(humanReadable, " minute", "m", 1)
	humanReadable = strings.Replace(humanReadable, " seconds", "s", 1)
	humanReadable = strings.Replace(humanReadable, " second", "s", 1)

	return humanReadable
}

func main() {
	var previousStatus status
	var nextStatus status

	log.Println("fetching previous data")
	err := previousStatus.fetchCurrent()
	if err != nil {
		log.Println(errors.Wrap(err, "error getting current status"))
		os.Exit(1)
	}

	if os.Getenv("REFRESH") == "" {
		log.Println("using previous data")
		nextStatus = previousStatus
	} else {
		log.Println("fetching latest data")
		nextStatus.fetchNew(previousStatus)
	}

	nextStatus.setCreatedAtStrings()

	if nextStatus == previousStatus {
		log.Println("no update required, exiting")
		os.Exit(0)
	}

	if os.Getenv("UPLOAD") == "" {
		log.Println("skipping upload")
		return
	}

	log.Println("formatting data for upload")
	statusJSON, err := json.Marshal(nextStatus)
	if err != nil {
		log.Fatal(errors.Wrap(err, "json generation error"))
		os.Exit(1)
	}

	log.Println("uploading data")
	err = upload(string(statusJSON))
	if err != nil {
		log.Fatal(errors.Wrap(err, "upload error"))
		os.Exit(1)
	}
	log.Println("completed successfully")
}
