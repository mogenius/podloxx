FROM golang:1.19-alpine AS builder

ENV CGO_ENABLED=1 GOOS=linux

RUN apk add --no-cache \
    libpcap-dev \
    g++ \
    perl-utils \
    curl \
    build-base \
    binutils-gold \
    bash \
    clang \
    llvm \
    libbpf-dev \
    linux-headers

ARG COMMIT_HASH=NOT_SET
ARG GIT_BRANCH=NOT_SET
ARG BUILD_TIMESTAMP=NOT_SET
ARG NEXT_VERSION=NOT_SET
ARG GITUSER
ARG GITPAT

RUN go env -w GOPRIVATE=github.com/mogenius
RUN apk add git
RUN git config --global url."https://$GITUSER:$GITPAT@github.com".insteadOf "https://github.com"


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-extldflags= \
  -X 'podloxx/version.GitCommitHash=${COMMIT_HASH}' \
  -X 'podloxx/version.Branch=${GIT_BRANCH}' \
  -X 'podloxx/version.BuildTimestamp=${BUILD_TIMESTAMP}' \
  -X 'podloxx/version.Ver=${NEXT_VERSION}'" -o bin/podloxx .


FROM alpine:latest
RUN apk add --no-cache \
    libpcap-dev bash

WORKDIR /app

COPY --from=builder ["/app/bin/podloxx", "."]
COPY --from=builder ["/app/.env", "/app/.env"]

ENV GIN_MODE=release

ENTRYPOINT [ "/app/podloxx", "cluster" ]