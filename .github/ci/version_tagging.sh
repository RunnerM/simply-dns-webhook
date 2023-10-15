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

git config --global user.email "ci-bot@pentek.dk"
git config --global user.name "runnerm-ci-bot"

git tag "$TAG"
git push origin "$TAG"


