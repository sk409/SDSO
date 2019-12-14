#!/bin/bash
GOOS=linux GOARCH=amd64 go build -o ./sdso ./client/app
cp ./sdso ./vulner/user/os/usr/bin/
docker-compose up -d --build