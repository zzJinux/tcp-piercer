#!/usr/bin/env bash
set -eu

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "$ROOT/utils/container_utils.sh"
eval "$("$ROOT/compose_env.sh" echoer)"


: ${COMPOSE}
: ${COMPOSE_FILE}
: ${COMPOSE_PROJECT_NAME}
: ${MODULE_PATH}
COMPOSE_P=$COMPOSE_PROJECT_NAME
SERVER_NAME=echoer
SERVER_PORT=8080


function setup() {
  $COMPOSE down --volumes
  $COMPOSE up -d
  router_contid=$($COMPOSE ps -q router)

  innersubnet_name=${COMPOSE_P}_router_inner
  outersubnet_name=${COMPOSE_P}_router_outer

  innersubnet=$(cont__get_network_subnet $innersubnet_name)
  outersubnet=$(cont__get_network_subnet $outersubnet_name)

  inner_router_ip=$(cont__get_ipaddr $router_contid $innersubnet_name)
  outer_router_ip=$(cont__get_ipaddr $router_contid $outersubnet_name)

  echoer_ip=$(cont__get_ipaddr $($COMPOSE ps -q $SERVER_NAME) ${COMPOSE_P}_router_outer)
  test_server=${echoer_ip%/*}:$SERVER_PORT

  # (on router) Configure the router to NAT packets from tp_client
  $COMPOSE exec -e SRC_SUBNET=$innersubnet -e DEST_SUBNET=$outersubnet -e TO_SOURCE=${outer_router_ip%/*} \
    router /scripts/setup_nat.sh 

  # (on tp_client) Add a new route from "tp_client" to "router"
  $COMPOSE exec tp_client ip route add $outersubnet via ${inner_router_ip%/*} dev eth0
}


function run() {
  $COMPOSE exec -e TEST_SERVER=$test_server tp_client /bin/sh -c \
    "/scripts/wait-for-command.sh -t 2 -c 'nc -z $test_server' && go test $MODULE_PATH/share/pnet"
}


function teardown() {
  $COMPOSE down --volumes
}


if [ "${ARTIFACT_DIR:-}" -a -d "${ARTIFACT_DIR}" ]; then
  envs_file="${ARTIFACT_DIR}/envs"
  if [ -e "$envs_file" ]; then
    source "$envs_file"
  fi
fi

# check if "test_server" has been set by TEST_CACHE
if ! [ ${test_server:+asdf} ]; then
  setup
fi

run
result=$?
if (exit $result); then
  teardown
fi

if [ "$envs_file" ]; then
  compose_shortcut="COMPOSE_FILE=${COMPOSE_FILE@Q} COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME@Q} $COMPOSE"
  declare -p compose_shortcut > "$envs_file"
  if ! (exit $result); then
    declare -p test_server > "$envs_file"
  fi
fi
