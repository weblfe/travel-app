#!/usr/bin/env bash

author=weblinuxgame
appName=travel-app

if [[ "${1}" != "" ]];then
  version="${1}"
else
  version=v3.16
fi

docker build -t ${author}/${appName}:${version} . && docker push ${author}/${appName}:${version}