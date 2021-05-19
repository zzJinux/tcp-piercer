#!/bin/sh
set -eu

# Required parameters
: ${IPTABLES:=iptables}
# SRC_SUBNET
# DEST_SUBNET
# TO_SOURCE

# inner_if=$(ip -br link | awk '$3 ~ /'$SRC_MACADDR'/ {print $1}')
# inner_if=${inner_if%@*}

rulespec="-s $SRC_SUBNET -d $DEST_SUBNET -j SNAT --to-source $TO_SOURCE"
if ! $IPTABLES -t nat -C POSTROUTING $rulespec 2>/dev/null; then
  $IPTABLES -t nat -A POSTROUTING $rulespec
fi
