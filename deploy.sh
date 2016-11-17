#!/bin/sh
commithash=$(git rev-parse --verify HEAD)
echo "building..."
GOOS=linux go build -ldflags "-X main.commithash=$commithash" .
echo "built."
echo "deploying..."
scp ./music_b austindoes.work:/home/austin/bin/
sudo scp -r ./static/* austindoes.work:/var/mbwww/
echo "donezo."
