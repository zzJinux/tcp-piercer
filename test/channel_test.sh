#!/usr/bin/env bash
set -eu

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "$ROOT/utils/container_utils.sh"

# This test usesd "echoer" setup
eval "$("$ROOT/compose_env.sh" server)"


: ${COMPOSE}
: ${COMPOSE_FILE}
: ${COMPOSE_PROJECT_NAME}
: ${MODULE_PATH}
COMPOSE_P=$COMPOSE_PROJECT_NAME
SERVER_NAME=tp_server
SERVER_PORT=9090


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

  server_ip=$(cont__get_ipaddr $(compose_fn ps -q $SERVER_NAME) ${COMPOSE_P}_router_outer)
  server_address=${server_ip%/*}:$SERVER_PORT

  # (on router) Configure the router to SNAT packets from tp_client
  compose_fn exec -e SRC_SUBNET=$innersubnet -e DEST_SUBNET=$outersubnet -e TO_SOURCE=${outer_router_ip%/*} \
    router /scripts/setup_nat.sh 

  # (on tp_client) Add a new route from "tp_client" to "router"
  compose_fn exec tp_client ip route add $outersubnet via ${inner_router_ip%/*} dev eth0
}


function run() {
  # execute both simulaneously, and then check if:
  #   1. both exits with zero
  #   2. compare res_client.json and res_server.json

  mkdir -p "$ROOT/_log"

  compose_fn exec tp_server /bin/sh -c \
    "go test $MODULE_PATH/test -run TestInitChannel -role=server -port=$SERVER_PORT -resultpath=/shared/res_server.json" \
    &>"$ROOT/_log/channel_test_server.txt" &
  _CLIENT_PID=$!

  # TODO: replace "sleep 3" with more reliable one
  compose_fn exec tp_client /bin/sh -c \
    "sleep 3 \
      && go test $MODULE_PATH/test -run TestInitChannel -role=client -serveraddr=$server_address -resultpath=/shared/res_client.json" \
    &>"$ROOT/_log/channel_test_client.txt" &
  _SERVER_PID=$!
  
  wait $_CLIENT_PID || { ret=$?; >&2 echo "client failed"; return $ret; }
  wait $_SERVER_PID || { ret=$?; >&2 echo "server failed"; return $ret; }

  copycon=$($DOCKER run --rm -d --mount type=volume,src="${COMPOSE_P}_shared-vol",dst=/shared,readonly alpine /bin/sleep 10)
  $DOCKER exec $copycon /bin/sh -c 'cat /shared/res_client.json' | jq -S . >"$ROOT/_log/res_client.json"
  $DOCKER exec $copycon /bin/sh -c 'cat /shared/res_server.json' | jq -S . >"$ROOT/_log/res_server.json"
  $DOCKER kill -s KILL $copycon
  
  diff "$ROOT/_log/res_client.json" "$ROOT/_log/res_server.json"
}



function teardown() {
  kill -0 $_SERVER_PID &>/dev/null && kill $_SERVER_PID
  kill -0 $_CLIENT_PID &>/dev/null && kill $_CLIENT_PID
  compose_fn down --volumes
}

setup
run
result=$?
if (exit $result); then
  teardown
fi
