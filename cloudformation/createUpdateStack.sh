#!/bin/bash

aws cloudformation create-stack --stack-name sujalEphemeralTweets --template-body file:///./EphemeralTweets.template