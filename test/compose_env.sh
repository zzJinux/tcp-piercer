#!/usr/bin/env bash
set -eu

if ! command -v yq &>/dev/null; then
  >&2 cat <<EOF
compose_env.sh: yq is required
See https://github.com/mikefarah/yq
EOF
  echo 'exit 1'
  exit 1
fi

composefile_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/composefiles" && pwd -P)"
base_composefile=${composefile_dir}/base.compose.yml
echoer_composefile=${composefile_dir}/echoer.compose.yml
server_composefile=${composefile_dir}/server.compose.yml

default_prefix=tp_test

export_and_print() {
  export "$@"
  declare -p "$@"
}

common() {
  : ${COMPOSE_PROJECT_NAME:=$default_prefix}
  : ${DOCKER:=docker}
  : ${COMPOSE:=docker-compose}
}

echoer() {
  COMPOSE_FILE="${base_composefile}:${echoer_composefile}"
  if [ "${GOMODCACHE:-}" ]; then
    read -r -d '' compose_file_stdin <<EOF || true
version: "3"
volumes:
  gomodcache:
    driver: local
    driver_opts:
      o: bind
      type: none
      device: $GOMODCACHE
services:
  tp_client:
    volumes:
      - gomodcache:/go/pkg/mod
EOF
    COMPOSE_FILE="${COMPOSE_FILE}:-"
  fi
}

server() {
  COMPOSE_FILE="${base_composefile}:${server_composefile}"
  if [ "${GOMODCACHE:-}" ]; then
    read -r -d '' compose_file_stdin <<EOF || true
version: "3"
volumes:
  gomodcache:
    driver: local
    driver_opts:
      o: bind
      type: none
      device: $GOMODCACHE
services:
  tp_client:
    volumes:
      - gomodcache:/go/pkg/mod
  tp_server:
    volumes:
      - gomodcache:/go/pkg/mod
EOF
    COMPOSE_FILE="${COMPOSE_FILE}:-"
  fi
}

common
case $1 in
echoer)
  echoer
  ;;
server)
  server
  ;;
esac

if [ "${compose_file_stdin:-}" ]; then
  declare -p compose_file_stdin
  compose_fn() {
    $COMPOSE "$@" <<<$compose_file_stdin
  }
else
  compose_fn() {
    $COMPOSE "$@"
  }
fi
declare -pf compose_fn
export_and_print COMPOSE_PROJECT_NAME DOCKER COMPOSE COMPOSE_FILE
