#!/usr/bin/env bash

dirNow=$(cd $(dirname "$0");pwd )

function stop() {
    # shellcheck disable=SC2009
    ps -ef | grep travel-app |grep -v grep | awk '{print $2}' | xargs kill -9
}

function update() {
    git pull && go build
}

function start() {
  nohup ./travel-app >app.log 2>&1 &
}

function main() {
    cd "${dirNow}"
    update
    stop
    start
}

main
