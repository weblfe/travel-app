#!/usr/bin/env bash

author=weblinuxgame
appName=travel
version=v3.16

docker build -t ${author}/${appName}:${version} . && docker push ${author}/${appName}:${version}