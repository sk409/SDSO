#!/bin/bash
go build -o ./sdso ./client/app
cp ./sdso ./vulner/user/os/usr/bin/