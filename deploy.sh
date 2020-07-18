#!/usr/bin/env bash

dirNow=$(dirname $(readlink -f "$0"))

function stop() {
    # shellcheck disable=SC2009
    ps -ef | grep travel-app |grep -v grep | awk '{print $2}' | xargs kill -9
}

function update() {
    git pull && go build
}

function start() {
  nohup "${dirNow}/travel-app" >app.log 2>&1 &
}

function main() {
    cd "${dirNow}"
    update
    stop
    start
}

main
