package main

import (
    "net/url"
    "os"
    "time"
    "encoding/json"
    "reflect"
    "strings"
    "strconv"

    "github.com/ChimeraCoder/anaconda"
    log "github.com/Sirupsen/logrus"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/secretsmanager"
)

var (
    logger            = log.New()

    // I want to set these during our main run after the logging is initialized
    maxTweetAge       = ""
    maxFavoriteAge 	  = ""
    secretARN         = ""
    ignoreTweets      = []string(nil)
    consumerKey       = ""
    consumerSecret    = ""
    accessToken       = ""
    accessTokenSecret = ""
)

func getignores(name string) []string {
    ignore_str := getenv(name)
    if len(ignore_str) > 0 {
        log.Info("There are tweets to ignore when deleting")
        return strings.Split(ignore_str, ",")
    }
    return nil
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

type Secret struct {
    TWITTER_CONSUMER_KEY string
    TWITTER_CONSUMER_SECRET string
    TWITTER_ACCESS_TOKEN string
    TWITTER_ACCESS_TOKEN_SECRET string
}

func getField(s *Secret, field string) string {
    r := reflect.ValueOf(s)
    f := reflect.Indirect(r).FieldByName(field)
    return f.String()
}

func getsecret(name string) string {
    sess, err := session.NewSession()
    if err != nil {
        log.Error("Got error creating new session")
        log.Error(err)
        os.Exit(1)
    }

    svc := secretsmanager.New(sess)
    secret, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: &secretARN})
    if err != nil {
        log.Error("Got error fetching secret")
        log.Error(err)
        os.Exit(1)
    }

    var s Secret
    err = json.Unmarshal([]byte(*secret.SecretString), &s)
    if err != nil {
        log.Error("Got error converting secret from JSON")
        log.Error(err)
        os.Exit(1)
    }
    secret_value := getField(&s, name)
    return secret_value
}

func getenv(name string) string {
    v := os.Getenv(name)
    if v == "" {
        log.Warn("Missing or empty requested environment variable " + name)
    }
    return v
}

func getTimeline(api *anaconda.TwitterApi, maxId int64) ([]anaconda.Tweet, error) {
    args := url.Values{}
    args.Add("count", "200")       // Twitter only returns most recent 20 tweets by default, so override
    args.Add("include_rts", "true") // When using count argument, RTs are excluded, so include them as recommended
    if (maxId > 0) {
        args.Add("max_id", strconv.FormatInt(maxId,10))
    }

    log.Info("Getting tweets: ", args)

    timeline, err := api.GetUserTimeline(args)
    if err != nil {
        return make([]anaconda.Tweet, 0), err
    }
    return timeline, nil
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func deleteFromTimeline(api *anaconda.TwitterApi, ageLimit time.Duration) {

    var maxTweetSeenId int64
    var loopCount int8

    for {

        timeline, err := getTimeline(api, maxTweetSeenId-1)

        log.Info("Timeline data length: ", len(timeline))

        // don't loop indefinitely and don't process if the timeline is empty. Probably no more tweets.
        timelineLength := len(timeline)
        if timelineLength == 0 || loopCount > 16 {
            break
        }

        if err != nil {
            log.Error("Could not get timeline")
        }
        for _, t := range timeline {
            createdTime, err := t.CreatedAtTime()
            if err != nil {
                log.Error("Couldn't parse time ", err)
            } else {
                if time.Since(createdTime) > ageLimit {
                    if stringInSlice(strconv.FormatInt(t.Id, 10), ignoreTweets) {
                        log.Info("Ignoring tweet while deleting " + strconv.FormatInt(t.Id, 10))
                        break
                    }
                    _, err := api.DeleteTweet(t.Id, true)
                    log.Info("DELETED: Age - ", time.Since(createdTime).Round(1*time.Minute), " - ", t.Text)
                    if err != nil {
                        log.Error("Failed to delete ", t.Id, "! ", err)
                    }
                }
            }

            maxTweetSeenId = t.Id

        }

        // short circuit if it's safely the "last" page. I'm just using the 190 number as an arbitrary cutoff.
        if timelineLength < 190 {
            break
        }

        loopCount++

    }

    log.Info("No more tweets to delete.")

}

func ephemeral() {
    // Initialize the loggin
    fmter := new(log.TextFormatter)
    fmter.FullTimestamp = true
    log.SetFormatter(fmter)
    log.SetLevel(log.InfoLevel)

    maxTweetAge       = getenv("MAX_TWEET_AGE")
    maxFavoriteAge    = getenv("MAX_FAVORITE_AGE")
    secretARN         = getenv("SECRET_ARN")

    ignoreTweets      = getignores("IGNORE_TWEETS")

    consumerKey       = getsecret("TWITTER_CONSUMER_KEY")
    consumerSecret    = getsecret("TWITTER_CONSUMER_SECRET")
    accessToken       = getsecret("TWITTER_ACCESS_TOKEN")
    accessTokenSecret = getsecret("TWITTER_ACCESS_TOKEN_SECRET")

    // We can continue without the max tweet age, but probably shouldn't start
    // deleting things since the user doesn't seem to know what they are doing
    if len(maxTweetAge) == 0 {
        log.Info("Max tweet age set to the default of 100 years")
        maxTweetAge   = "100y"
    }
    ht, _ := time.ParseDuration(maxTweetAge)

    if len(maxFavoriteAge) == 0 {
    	log.Info("Max favorite age set to the default of 100 years")
    	maxFavoriteAge 	= "100y"
    }
    hf, _ := time.ParseDuration(maxFavoriteAge)


    // We can't continue without the secretARN
    if len(secretARN) == 0 {
        log.Error("We can't continue without the secret ARN")
        os.Exit(1)
    }

    anaconda.SetConsumerKey(consumerKey)
    anaconda.SetConsumerSecret(consumerSecret)
    api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
    api.SetLogger(anaconda.BasicLogger)

    deleteFromTimeline(api, ht)

    unfavorite(api, hf)
}

func main() {

    lambda.Start(ephemeral)

}
