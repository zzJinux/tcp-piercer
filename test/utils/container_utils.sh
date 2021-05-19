#!/usr/bin/env bash

# $1 - network name
cont__get_network_subnet() {
  $DOCKER network inspect -f '{{ (index .IPAM.Config 0).Subnet }}' $1
  # outputs in CIDR format
}

# $1 - container name
# $2 - environment variable name
cont__read_env() {
  $DOCKER exec $1 /bin/sh -c "echo \${$2}"
}

# $1 - container name
# $2 - network name
cont__get_ipaddr() {
  $DOCKER container inspect -f "{{.NetworkSettings.Networks.$2.IPAddress}}" $1
}

# $1 - container name
# $2 - network name
cont__get_macaddr() {
  $DOCKER container inspect -f "{{.NetworkSettings.Networks.$2.MacAddress}}" $1
}
