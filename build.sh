#!/usr/bin/env bash

author=weblinuxgame
app-name=travel
version=v3.16

docker build -t ${author}/travel:${version} . && docker push ${author}/travel:${version}