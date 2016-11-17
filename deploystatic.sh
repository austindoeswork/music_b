#!/bin/sh
echo "deploying static..."
sudo scp -r ./static/* austindoes.work:/var/mbwww/
echo "donezo."
