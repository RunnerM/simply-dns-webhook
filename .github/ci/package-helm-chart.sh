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

git config --global user.email "ci-bot@runnerm.com"
git config --global user.name "simple-dns-webhook CI robot"

git add --all
git commit -m "Chore: Update helm chart for version $VERSION"
git push