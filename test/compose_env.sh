#!/usr/bin/env bash
set -eu

composefile_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/composefiles" && pwd -P)"
base_composefile=${composefile_dir}/base.compose.yml
echoer_composefile=${composefile_dir}/echoer.compose.yml
server_composefile=${composefile_dir}/server.compose.yml

default_prefix=tp_test

export_and_print() {
  export "$@"
  declare -p "$@"
  echo
}

common() {
  : ${COMPOSE_PROJECT_NAME:=$default_prefix}
  : ${DOCKER:=docker}
  : ${COMPOSE:=docker-compose}
  export_and_print COMPOSE_PROJECT_NAME DOCKER COMPOSE
}

echoer() {
  : ${COMPOSE_FILE:=${base_composefile}:${echoer_composefile}}
  export_and_print COMPOSE_FILE
}

case $1 in
echoer)
  common
  echoer
  ;;
esac
