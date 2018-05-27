package main

import (
	"net/url"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	consumerKey       = getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")
	maxTweetAge       = getenv("MAX_TWEET_AGE")
	maxFavoriteAge    = getenv("MAX_FAVORITE_AGE")
	logger            = log.New()
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

func getTimeline(api *anaconda.TwitterApi) ([]anaconda.Tweet, error) {
	args := url.Values{}
	args.Add("count", "200")       // Twitter only returns most recent 20 tweets by default, so override
	args.Add("include_rts", "true") // When using count argument, RTs are excluded, so include them as recommended
	timeline, err := api.GetUserTimeline(args)
	if err != nil {
		return make([]anaconda.Tweet, 0), err
	}
	return timeline, nil
}

func deleteFromTimeline(api *anaconda.TwitterApi, ageLimit time.Duration) {
	timeline, err := getTimeline(api)

	if err != nil {
		log.Error("Could not get timeline")
	}
	for _, t := range timeline {
		createdTime, err := t.CreatedAtTime()
		if err != nil {
			log.Error("Couldn't parse time ", err)
		} else {
			if time.Since(createdTime) > ageLimit {
				_, err := api.DeleteTweet(t.Id, true)
				log.Info("DELETED: Age - ", time.Since(createdTime).Round(1*time.Minute), " - ", t.Text)
				if err != nil {
					log.Error("Failed to delete! ", err)
				}
			}
		}
	}
	log.Info("No more tweets to delete.")
}

func getFavorites(api *anaconda.TwitterApi) ([]anaconda.Tweet, error) {
	args := url.Values{}
	args.Add("count", "200")       // Twitter only returns most recent 20 tweets by default, so override
	timeline, err := api.GetFavorites(args)
	if err != nil {
		return make([]anaconda.Tweet, 0), err
	}
	return timeline, nil
}

func unfavorite(api *anaconda.TwitterApi, ageLimit time.Duration) {
	favorites, err := getFavorites(api)

	if err != nil {
		log.Error("Could not get favorites")
	}
	for _, f := range favorites {
		createdTime, err := f.CreatedAtTime()
		if err != nil {
			log.Error("Couldn't parse time ", err)
		} else {
			if time.Since(createdTime) > ageLimit {
				_, err := api.Unfavorite(f.Id)
				log.Info("UNFAVORITED: Age - ", time.Since(createdTime).Round(1*time.Minute), " - ", f.Text)
				if err != nil {
					log.Error("Failed to unfavorite! ", err)
				}
			}
		}
	}
	log.Info("No more tweets to unfavorite.")
}

func ephemeral() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	api.SetLogger(anaconda.BasicLogger)

	fmter := new(log.TextFormatter)
	fmter.FullTimestamp = true
	log.SetFormatter(fmter)
	log.SetLevel(log.InfoLevel)

	ht, _ := time.ParseDuration(maxTweetAge)
	deleteFromTimeline(api, ht)

	hf, _ := time.ParseDuration(maxFavoriteAge)
	unfavorite(api, hf)
}

func main() {
	lambda.Start(ephemeral)
}
