#!/bin/sh
#This is a script for building and indexing the chart.

VERSION=$1
if [ -z "$VERSION" ]; then
  echo "No version supplied"
  exit 1
fi

IMAGE_TAG="v$VERSION"


cd deploy
helm lint simply-dns-webhook
helm package simply-dns-webhook --version $VERSION --app-version $VERSION
helm repo index . --url https://runnerm.github.io/simply-dns-webhook/