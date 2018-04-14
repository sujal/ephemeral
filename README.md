# ephemeral: automatically delete your old Tweets with AWS Lambda

**ephemeral** is a Twitter timeline grooming program that runs for pretty much free on AWS Lambda. The code is forked from Adam Drake's excellent [Harold](https://github.com/adamdrake/harold) butler-like bot and refactored for Lambda.

You can use ephemeral to automatically delete all tweets from your timeline that are older than a certain number of hours that you can choose. For instance, you can ensure that your tweets are deleted after one week (168h), or one day (24h).

The program will run once for each execution based on the trigger/schedule you set in AWS Lambda.

# Twitter API

You will need to [create a new Twitter application and generate API keys](https://apps.twitter.com/). The program assumes the following secrets are set under a single secret in [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/).  
Secrets Manager does not have a Free Tier option but is not expensive for the added security.

```
TWITTER_CONSUMER_KEY
TWITTER_CONSUMER_SECRET
TWITTER_ACCESS_TOKEN
TWITTER_ACCESS_TOKEN_SECRET
```

And the following needs to be set as environment variables.
```
MAX_TWEET_AGE
SECRET_ARN
```

`MAX_TWEET_AGE` expects a value of hours, example: `MAX_TWEET_AGE=72h` (Make sure to end the value with `h` or equivalent)  
`SECRET_ARN` expects the full ARN of the secret, example: `arn:aws:secretsmanager:us-east-2:000000000000:secret:ephemeralTweets-Uf8NON`

You can set these variables in AWS Lambda when you create your Lambda function. For a full walkthrough with screenshots on creating a Lambda function and uploading the code, read [this blog post](https://vickylai.com/verbose/aws-lambda/). Skip to setting environment variables at [this link](https://vickylai.com/verbose/aws-lambda/#2-configure-your-function).

# update.sh

This handy bash script is included to help you upload your function code to Lambda. It requires [AWS Command Line Interface](https://aws.amazon.com/cli/). To set up, do `pip install awscli` and follow these instructions for [Quick Configuration](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html).
