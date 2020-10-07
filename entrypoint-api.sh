#!/usr/bin/env sh

# shellcheck disable=SC2112
function updateHost() {
    host=${1}
    echo "${host} $(cat /etc/hosts.d/${hosts} | grep ${hosts} | awk '{print $1}')" >> /etc/hosts
}

# 启动
function start() {
    # shellcheck disable=SC2039
    if [[ ! -e /etc/hosts.d/lock ]];then
      updateHost redis-server
      updateHost mysql-server
      updateHost mongo-server
      updateHost nats-server
      touch /etc/hosts.d/lock
    fi
    /data/www/app/api-server
}

start