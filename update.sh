#!/bin/bash

git pull

export GO111MODULE=off

go build -o excel2config cmd/main.go

killall excel2config

sleep 2

nohup ./excel2config -conf ./configs/ 2>&1 &

echo "restart succeed"