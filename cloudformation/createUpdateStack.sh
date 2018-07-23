#!/bin/bash

TEMPLATE_URL=$1
SECRET_ARN=$2
MAX_TWEET_AGE=$4
MAX_FAVORITE_AGE=$3
IGNORE_TWEET_LIST=$5

if [ -z $TEMPLATE_URL ]; then
	echo "Template URL is required"
	exit 1
fi

if [ -z $SECRET_ARN ]; then
	echo "Secret ARN is required"
	exit 1
fi

if [ -z $MAX_TWEET_AGE ]; then
	echo "Max Tweet Age is required"
	exit 1
fi

if [ -z $MAX_FAVORITE_AGE ]; then
	echo "Max Favorite Age is required"
	exit 1
fi


aws cloudformation create-stack --stack-name sujalEphemeralTweets --template-url $TEMPLATE_URL --capabilities CAPABILITY_IAM --parameters ParameterKey=SecretArn,ParameterValue=$SECRET_ARN ParameterKey=MaxTweetAge,ParameterValue=$MAX_TWEET_AGE ParameterKey=MaxFavoriteAge,ParameterValue=$MAX_FAVORITE_AGE ParameterKey=IgnoreTweets,ParameterValue=$IGNORE_TWEET_LIST
