#!/bin/bash
set -e

GIT_BRANCH=$(git branch | grep \* | cut -d ' ' -f2 | tr '[:upper:]' '[:lower:]')
NEXT_VERSION=1.0.2
COMMIT_HASH=$(git rev-parse --short HEAD)
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
BUILD_TIMESTAMP=$(date)
CGO_ENABLED=1

GOOS=darwin
GOARCH=amd64
go build -ldflags="-extldflags= \
  -X 'podloxx-collector/version.GitCommitHash=${COMMIT_HASH}' \
  -X 'podloxx-collector/version.Branch=${GIT_BRANCH}' \
  -X 'podloxx-collector/version.BuildTimestamp=${BUILD_TIMESTAMP}' \
  -X 'podloxx-collector/version.Ver=${NEXT_VERSION}'" -o bin/podloxx-${NEXT_VERSION}-$GOOS-$GOARCH .

GOOS=darwin 
GOARCH=arm64
go build -ldflags="-extldflags= \
  -X 'podloxx-collector/version.GitCommitHash=${COMMIT_HASH}' \
  -X 'podloxx-collector/version.Branch=${GIT_BRANCH}' \
  -X 'podloxx-collector/version.BuildTimestamp=${BUILD_TIMESTAMP}' \
  -X 'podloxx-collector/version.Ver=${NEXT_VERSION}'" -o bin/podloxx-${NEXT_VERSION}-$GOOS-$GOARCH .