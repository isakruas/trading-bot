#!/usr/bin/env bash
set -e

# Gets the latest semantic tag
last_tag=$(git tag --sort=-v:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | head -n 1)

# If there is no previous tag
if [ -z "$last_tag" ]; then
  echo "0.0.1"
  exit 0
fi

version=${last_tag#v}
major=$(echo "$version" | cut -d. -f1)
minor=$(echo "$version" | cut -d. -f2)
patch=$(echo "$version" | cut -d. -f3)

commits=$(git log "$last_tag"..HEAD --pretty=format:"%s%n%b")

# Checks for commit messages that indicate version bumps
if echo "$commits" | grep -q "BREAKING CHANGE:"; then
  major=$((major + 1))
  minor=0
  patch=0
elif echo "$commits" | grep -q "^feat"; then
  minor=$((minor + 1))
  patch=0
else
  patch=$((patch + 1))
fi

# Outputs the new semantic version
echo "$major.$minor.$patch"
