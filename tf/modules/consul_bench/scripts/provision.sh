#!/bin/bash
set -ex

chmod +x ./consul-bench

while true; do
    curl $1:8500/v1/agent/self > /dev/null && break
    sleep 1
done

echo $1 > target

echo "./consul-bench -rpc-addr $1:8300 -consul $1:8500 -time 30s -register 10 > results.log" > start
chmod +x start
