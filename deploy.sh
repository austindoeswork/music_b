#!/bin/sh
GOOS=linux go build .
scp ./music_b austindoes.work:/home/austin/bin/
sudo scp -r ./static/* austindoes.work:/var/mbwww/
echo "donezo."
