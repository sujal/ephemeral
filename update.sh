#!/bin/bash

if [ -z $AWS_REGION ]; then
	AWS_REGION=${AWS_DEFAULT_REGION:-us-east-2}
	echo "defaulting to $AWS_REGION"
fi

which go >/dev/null 2>&1
if [ $? -eq 0 ]; then
    go get -d ./.
    go build -o ephemeralTweets && \
        zip /tmp/ephemeralTweets.zip ephemeralTweets && \
        aws --region $AWS_REGION lambda update-function-code --function-name ephemeralTweets --zip-file fileb:///tmp/ephemeralTweets.zip && \
        rm -rf /tmp/ephemeralTweets.zip ./ephemeralTweets
else
    echo "You need to have golang installed"
    exit 1
fi
