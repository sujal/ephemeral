# ephemeral: automatically delete your old Tweets with AWS Lambda

----

**NOTE**: This is forked from [Vicky Lai's project](https://github.com/vickylai/ephemeral) with work from several other authors merged in. Specifically:

* the [fork by osvaldsson](https://github.com/osvaldsson/ephemeral) - adds AWS Secrets Manager support and misc script improvements
* the [fork by grantwinney](https://github.com/grantwinney/ephemeral) - adds the ability to delete favorites older than a time window.

The major addition from my perspective has been the Cloudformation templates. I'm still getting more comfortable with CF - I'm a relative novice, so feedback is welcome.

See the end of this readme for CloudFormation instructions.

----


**ephemeral** is a Twitter timeline grooming program that runs for pretty much free on AWS Lambda. The code is forked from Adam Drake's excellent [Harold](https://github.com/adamdrake/harold) butler-like bot and refactored for Lambda.

You can use ephemeral to automatically delete all tweets and favorites from your timeline that are older than a certain number of hours that you can choose. For instance, you can ensure that your tweets are deleted after one week (168h), or one day (24h).

The program will run once for each execution based on the trigger/schedule you set in AWS Lambda. It will delete up to 200 expired tweets (per-request limit set by Twitter's API) each run.

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
MAX_FAVORITE_AGE
SECRET_ARN
IGNORE_TWEETS (optional)
```

`MAX_TWEET_AGE` and `MAX_FAVORITE_AGE` expects a value of hours, example: `MAX_TWEET_AGE=72h` (Make sure to end the value with `h` or equivalent)
`SECRET_ARN` expects the full ARN of the secret, example: `arn:aws:secretsmanager:us-east-2:000000000000:secret:ephemeralTweets-Uf8NON`
`IGNORE_TWEETS` expects a comma delimited string of tweet IDs that you do not want to delete.

You can set these variables in AWS Lambda when you create your Lambda function. For a full walkthrough with screenshots on creating a Lambda function and uploading the code, read [this blog post](https://vickylai.com/verbose/free-twitter-bot-aws-lambda/). Skip to setting environment variables at [this link](https://vickylai.com/verbose/free-twitter-bot-aws-lambda/#2-configure-your-function).

# update.sh

This handy bash script is included to help you upload your function code to Lambda. It requires [AWS Command Line Interface](https://aws.amazon.com/cli/). To set up, do `pip install awscli` and follow these instructions for [Quick Configuration](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html).

# CloudFormation

This entire project is best suited for people with AWS experience of some kind. I will work on better instructions soon, but for now, this is geared toward people familiar with AWS.

There is a Cloudformation template included in the project, along with a few scripts to take advantage of them. You can, however, follow these steps to quickly get started:

1. Add the 4 Twitter secrets to [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/). Digital Ocean has a [good list of steps to follow](https://www.digitalocean.com/community/tutorials/how-to-create-a-twitter-app) when creating your app. Name your app with a unique name and generate all of the credentials (basically everything up to and including Step 2).
2. Take note of the ARN for the new secret you've created.
3. Tap the Launch Stack button to get started:

<a href="https://console.aws.amazon.com/cloudformation/home?region=us-east-1#/stacks/new?stackName=EphemeralTweets&templateURL=https://code.eb.forchesoftware.com/code/EphemeralTweets.template">![launch stack](https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png)</a>

That should get you going. You'll need to enter in the times you want for the age parameters (See the section above for details). You will also need the ARN for the secreate you created in step 1.
