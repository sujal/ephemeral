package main

import (
    "net/url"
    "os"
    "time"
    "fmt"
    "encoding/json"
    "reflect"

    "github.com/ChimeraCoder/anaconda"
    log "github.com/Sirupsen/logrus"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/secretsmanager"
)

var (
    logger            = log.New()

    maxTweetAge       = getenv("MAX_TWEET_AGE")
    secretARN         = getenv("SECRET_ARN")

    consumerKey       = getsecret("TWITTER_CONSUMER_KEY")
    consumerSecret    = getsecret("TWITTER_CONSUMER_SECRET")
    accessToken       = getsecret("TWITTER_ACCESS_TOKEN")
    accessTokenSecret = getsecret("TWITTER_ACCESS_TOKEN_SECRET")
)

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
        fmt.Println("Got error creating new session")
        fmt.Println(err.Error())
        os.Exit(1)
    }

    svc := secretsmanager.New(sess)
    secret, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: &secretARN})
    if err != nil {
        fmt.Println("Got error fetching secret")
        fmt.Println(err.Error())
        os.Exit(1)
    }

    var s Secret
    err = json.Unmarshal([]byte(*secret.SecretString), &s)
    if err != nil {
        fmt.Println("Got error converting secret from JSON")
        fmt.Println(err.Error())
        os.Exit(1)
    }
    secret_value := getField(&s, name)
    return secret_value
}

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

func ephemeral() {
    anaconda.SetConsumerKey(consumerKey)
    anaconda.SetConsumerSecret(consumerSecret)
    api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
    api.SetLogger(anaconda.BasicLogger)

    fmter := new(log.TextFormatter)
    fmter.FullTimestamp = true
    log.SetFormatter(fmter)
    log.SetLevel(log.InfoLevel)

    h, _ := time.ParseDuration(maxTweetAge)

    deleteFromTimeline(api, h)

}

func main() {

    lambda.Start(ephemeral)

}
