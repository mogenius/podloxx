# Thanks to Stefan Buck. We took some code from his snipped: 
# License: MIT https://gist.github.com/stefanbuck/ce788fee19ab6eb0b4447a85fc99f447 Author: Stefan Buck
#
# Example ./build-client.sh NEXT_VERSION=1.0.1 GOOS=darwin GOARCH=amd64 github_api_token=XXX
#

#!/bin/bash
set -ex

xargs=$(which gxargs || which xargs)

# Validate settings.
[ "$TRACE" ] && set -x

CONFIG=$@

for line in $CONFIG; do
  eval "$line"
done

echo $NEXT_VERSION

if [[ -z "${NEXT_VERSION}" ]]; then
  echo "NEXT_VERSION is undefined."; exit 0
fi

if [[ -z "${GOOS}" ]]; then
  echo "GOOS is undefined."; exit 0
fi

if [[ -z "${GOARCH}" ]]; then
  echo "GOARCH is undefined."; exit 0
fi

if [[ -z "${github_api_token}" ]]; then
  echo "github_api_token is undefined."; exit 0
fi

owner=mogenius
repo=podloxx
GIT_BRANCH=$(git branch | grep \* | cut -d ' ' -f2 | tr '[:upper:]' '[:lower:]')
COMMIT_HASH=$(git rev-parse --short HEAD)
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
BUILD_TIMESTAMP=$(date)
CGO_ENABLED=1
tag="v$NEXT_VERSION"
filename="bin/podloxx-${NEXT_VERSION}-$GOOS-$GOARCH"

go mod download

go build -ldflags="-extldflags= \
  -s -w \
  -X 'podloxx-collector/version.GitCommitHash=${COMMIT_HASH}' \
  -X 'podloxx-collector/version.Branch=${GIT_BRANCH}' \
  -X 'podloxx-collector/version.BuildTimestamp=${BUILD_TIMESTAMP}' \
  -X 'podloxx-collector/version.Ver=${NEXT_VERSION}'" -o $filename .

# Define variables.
GH_API="https://api.github.com"
GH_REPO="$GH_API/repos/$owner/$repo"
GH_TAGS="$GH_REPO/releases/tags/$tag"
AUTH="Authorization: token $github_api_token"
WGET_ARGS="--content-disposition --auth-no-challenge --no-cookie"
CURL_ARGS="-LJO#"

if [[ "$tag" == 'LATEST' ]]; then
  GH_TAGS="$GH_REPO/releases/latest"
fi

# Validate token.
curl -o /dev/null -sH "$AUTH" $GH_REPO || { echo "Error: Invalid repo, token or network issue!";  exit 1; }

# Create the release
GH_REL="https://api.github.com/repos/$owner/$repo/releases"
curl -X POST "$GITHUB_OAUTH_BASIC" -H "Authorization: token $github_api_token" -H "Content-Type: application/octet-stream" -d "{\"tag_name\":\"$tag\",\"target_commitish\":\"develop\",\"name\":\"$tag\",\"body\":\"Description of the release\",\"draft\":false,\"prerelease\":false,\"generate_release_notes\":false}" $GH_REL

# Read asset tags.
response=$(curl -sH "$AUTH" $GH_TAGS)

# Get ID of the asset based on given filename.
eval $(echo "$response" | grep -m 1 "id.:" | grep -w id | tr : = | tr -cd '[[:alnum:]]=')
[ "$id" ] || { echo "Error: Failed to get release id for tag: $tag"; echo "$response" | awk 'length($0)<100' >&2; exit 1; }

# Upload asset
echo "Uploading asset... "

# Construct url
GH_ASSET="https://uploads.github.com/repos/$owner/$repo/releases/$id/assets?name=$(basename $filename)"

curl "$GITHUB_OAUTH_BASIC" --data-binary @"$filename" -H "Authorization: token $github_api_token" -H "Content-Type: application/octet-stream" $GH_ASSET