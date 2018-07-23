#!/bin/bash

FUNCTION_NAME=$1

if [ -z $FUNCTION_NAME ]; then
	echo "Function name is a required argument."
	exit 1
fi

if [ -z $AWS_REGION ]; then
	AWS_REGION=${AWS_DEFAULT_REGION:-us-east-2}
	echo "defaulting to $AWS_REGION"
fi

which go >/dev/null 2>&1
if [ $? -eq 0 ]; then
    go get -d ./.
    GOOS=linux GOARCH=amd64 go build -o ephemeralTweets && \
        zip /tmp/ephemeralTweets.zip ephemeralTweets && \
        aws s3 cp /tmp/ephemeralTweets.zip s3://code.eb.forchesoftware.com/code/ephemeralTweets.zip && \
        aws --region $AWS_REGION lambda update-function-code --function-name $FUNCTION_NAME --s3-bucket code.eb.forchesoftware.com --s3-key code/ephemeralTweets.zip && \
        rm -rf /tmp/ephemeralTweets.zip ./ephemeralTweets
else
    echo "You need to have golang installed"
    exit 1
fi
