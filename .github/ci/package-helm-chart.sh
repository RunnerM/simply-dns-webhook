#!/bin/sh
#This is a script for building and indexing the chart.

VERSION=$1
if [ -z "$VERSION" ]; then
  echo "No version supplied"
  exit 1
fi

IMAGE_TAG="v$VERSION"


cd deploy/simply-dns-webhook

sed -i '/tag:/c\  tag: $IMAGE_TAG' values.yaml
sed -i '/version: /c\version: $IMAGE_TAG' Chart.yaml
sed -i '/appVersion: "/c\appVersion: "$IMAGE_TAG"' Chart.yaml

cd ..
helm lint simply-dns-webhook
helm package simply-dns-webhook
helm repo index . --url https://runnerm.github.io/simply-dns-webhook/

git config --global user.email "ci-bot@runnerm.com"
git config --global user.name "simple-dns-webhook CI robot"

git add --all
git commit -m "Chore: Update helm chart for version $VERSION"
git push
