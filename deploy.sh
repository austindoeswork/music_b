#!/bin/sh
commithash=$(git rev-parse --verify HEAD)
echo "building..."
GOOS=linux go build -ldflags "-X main.commithash=$commithash" .
echo "built."
echo "deploying..."
scp ./music_b 138.197.26.172:/home/austin/bin/
scp -r ./static/* 138.197.26.172:/var/mbwww/
rm ./music_b
echo "donezo."
