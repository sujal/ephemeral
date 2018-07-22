#!/bin/bash

aws cloudformation create-stack --stack-name sujalEphemeralTweets --template-url s3://code.eb.forchesoftware.com/templates/EphemeralTweets.template