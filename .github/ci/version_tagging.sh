#!/bin/sh

VERSION=$1

if [ -z "$VERSION" ]; then
  echo "No version supplied"
  exit 1
fi

TAG="v$VERSION"

if git rev-parse "$TAG" >/dev/null 2>&1; then
  echo "Tag $TAG already exists"
  exit 1
fi

git config --global user.email "ci-bot@runnerm.com"
git config --global user.name "simple-dns-webhook CI robot"

git tag "$TAG"
git push origin "$TAG"


