#!/usr/bin/env bash
set -eu
shopt -s nullglob

FLUSH_ARTIFACTS=
case $1 in
  -f) FLUSH_ARTIFACTS=.; shift
esac

test_artifact_dir="test/artifacts/$1"
mkdir -p "$test_artifact_dir"
if [ $FLUSH_ARTIFACTS ]; then
  rm -f "$test_artifact_dir"/*
fi

ARTIFACT_DIR=$test_artifact_dir "test/${1}.sh" || true
