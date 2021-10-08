#!/usr/bin/env bash
set -eu

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "$ROOT/utils/container_utils.sh"

# This test usesd "echoer" setup
eval "$("$ROOT/compose_env.sh" echoer)"


: ${COMPOSE}
: ${COMPOSE_FILE}
: ${COMPOSE_PROJECT_NAME}
: ${MODULE_PATH}
COMPOSE_P=$COMPOSE_PROJECT_NAME
SERVER_NAME=echoer
SERVER_PORT=8080


function setup() {
  compose_fn down --volumes
  compose_fn up -d
  router_contid=$(compose_fn ps -q router)

  innersubnet_name=${COMPOSE_P}_router_inner
  outersubnet_name=${COMPOSE_P}_router_outer

  innersubnet=$(cont__get_network_subnet $innersubnet_name)
  outersubnet=$(cont__get_network_subnet $outersubnet_name)

  inner_router_ip=$(cont__get_ipaddr $router_contid $innersubnet_name)
  outer_router_ip=$(cont__get_ipaddr $router_contid $outersubnet_name)

  echoer_ip=$(cont__get_ipaddr $(compose_fn ps -q $SERVER_NAME) ${COMPOSE_P}_router_outer)
  echo_server=${echoer_ip%/*}:$SERVER_PORT

  # (on router) Configure the router to NAT packets from tp_client
  compose_fn exec -e SRC_SUBNET=$innersubnet -e DEST_SUBNET=$outersubnet -e TO_SOURCE=${outer_router_ip%/*} \
    router /scripts/setup_nat.sh 

  # (on tp_client) Add a new route from "tp_client" to "router"
  compose_fn exec tp_client ip route add $outersubnet via ${inner_router_ip%/*} dev eth0
}


function run() {
  compose_fn exec -e ECHO_SERVER=$echo_server tp_client /bin/sh -c \
    "/scripts/wait-for-command.sh -t 2 -c 'nc -z $echo_server' && go test $MODULE_PATH/share/pnet"
}


function teardown() {
  compose_fn down --volumes
}

setup
run
result=$?
if (exit $result); then
  teardown
fi
