#!/usr/bin/env bash
# Sourced by both builders. Only an explicit matching release version is stable.
community_version="$(tr -d '\r\n' < version/VERSION)"
if [[ ! "$community_version" =~ ^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]; then
  echo "Invalid version/VERSION: expected MAJOR.MINOR.PATCH" >&2
  exit 1
fi
ARTEX_BUILD_VERSION="${ARTEX_BUILD_VERSION:-${community_version}-dev}"
ARTEX_BUILD_VERSION="${ARTEX_BUILD_VERSION#v}"
if [[ "$ARTEX_BUILD_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ && "$ARTEX_BUILD_VERSION" != "$community_version" ]]; then
  echo "Release version must match version/VERSION" >&2
  exit 1
fi
