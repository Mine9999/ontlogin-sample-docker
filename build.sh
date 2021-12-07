#/bin/bash

docker run --rm -w /go/src/app -v $PWD/backend:/go/src/app golang:1.17 bash -c "go build"
docker run --rm -w /opt/app -v $PWD/frontend:/opt/app node:lts-stretch bash -c "npm install ; npm run build"