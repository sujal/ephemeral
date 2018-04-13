!#/bin/bash

go build && \
zip /tmp/ephemeral.zip ephemeral && \
aws --region us-east-2 lambda update-function-code --function-name ephemeralTweets --zip-file fileb:///tmp/ephemeral.zip && \
rm /tmp/ephemeral.zip
