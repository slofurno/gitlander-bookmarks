#!/bin/sh

a0="tcp://127.0.0.1:40890"
a1="tcp://127.0.0.1:40891"
a2="tcp://127.0.0.1:40892"
a3="tcp://127.0.0.1:40893"

pids=""

./cluster $a0 $a1 $a2 & pids="$pids $!"
./cluster $a1 $a2 $a3 & pids="$pids $!"
./cluster $a2 $a3 & pids="$pids $!"
./cluster $a3 $a0 & pids="$pids $!"

echo $pids

sleep 6 

kill $pids


