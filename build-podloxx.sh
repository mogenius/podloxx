#!/bin/bash

dir=$(dirname -- "$( readlink -f -- "$0"; )";)

cd ${dir}/ui

npm install -g @angular/cli
npm --force i
npm run build

cd ${dir}/

echo "Build podloxx"
go mod download

go build -ldflags="-extldflags= -s -w" -o ${dir}/bin/podloxx .

echo "Run podloxx: "
echo "${dir}/bin/podloxx start"
