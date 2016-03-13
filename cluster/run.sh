#!/bin/sh

a0="tcp://127.0.0.1:40890"
a1="tcp://127.0.0.1:40891"
a2="tcp://127.0.0.1:40892"
a3="tcp://127.0.0.1:40893"

pids=""

./cluster cluster1 :11411 tcp://:11400 $a0 $a1 $a2 & pids="$pids $!"
./cluster cluster2 :11412 tcp://:11402 $a1 $a2 $a3 & pids="$pids $!"
./cluster cluster3 :11413 tcp://:11403 $a2 $a3 & pids="$pids $!"
./cluster cluster4 :11414 tcp://:11404 $a3 $a0 & pids="$pids $!"

echo $pids

sleep 30

kill $pids


