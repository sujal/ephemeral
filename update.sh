#!/bin/bash

which go >/dev/null 2>&1
if [ $? -eq 0 ]; then
    go get -d ./.
    go build -o ephemeralTweets && \
        zip /tmp/ephemeralTweets.zip ephemeralTweets && \
        aws --region us-east-2 lambda update-function-code --function-name ephemeralTweets --zip-file fileb:///tmp/ephemeralTweets.zip && \
        rm -rf /tmp/ephemeralTweets.zip ./ephemeralTweets
else
    echo "You need to have golang installed"
    exit 1
fi
