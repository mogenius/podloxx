#!/bin/bash
set -e

CR_HOST=ghcr.io
OWNER=mogenius
NAME=podloxx
GIT_BRANCH=$(git branch | grep \* | cut -d ' ' -f2 | tr '[:upper:]' '[:lower:]')
NEXT_VERSION=1.0.4
COMMIT_HASH=$(git rev-parse --short HEAD)
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
BUILD_TIMESTAMP=$(date)

echo "### LOGIN ###"
if [ "$CR_PAT" = '' ]
then
  echo "Please export CR_PAT and try again. (https://github.com/settings/tokens -> Permissions: packages:read,write,delete)"
  exit 1
fi
echo $CR_PAT | docker login ghcr.io -u USERNAME --password-stdin

if [ "$GITHUBUSER" = '' ]
then
  echo "Please export GITHUBUSER and try again."
  exit 1
fi
if [ "$GITHUBPAT" = '' ]
then
  echo "Please export GITHUBPAT and try again. (https://github.com/settings/tokens -> Permissions: repos:ALL)"
  exit 1
fi

DOCKER_REPO=$CR_HOST/$OWNER/$NAME

DOCKER_TAGGED_BUILDS=("$DOCKER_REPO:latest" "$DOCKER_REPO:$NEXT_VERSION")

echo "### BUILD CONTAINER $NEXT_VERSION ### ${DOCKER_TAGGED_BUILDS[@]}"
DOCKER_TAGS_ARGS=$(echo ${DOCKER_TAGGED_BUILDS[@]/#/-t }) # "-t FIRST_TAG -t SECOND_TAG ..."
DOCKER_BUILDKIT=1 docker build --platform linux/amd64 $DOCKER_TAGS_ARGS --build-arg NEXT_VERSION="$NEXT_VERSION" --build-arg BUILD_TIMESTAMP="$BUILD_TIMESTAMP" --build-arg GIT_BRANCH="$GIT_BRANCH" --build-arg COMMIT_HASH="$COMMIT_HASH" --build-arg GITUSER="$GITHUBUSER" --build-arg GITPAT="$GITHUBPAT" -f Dockerfile .

for DOCKER_TAG in "${DOCKER_TAGGED_BUILDS[@]}"
do
  echo pushing "$DOCKER_TAG"
  docker push "$DOCKER_TAG"
done
