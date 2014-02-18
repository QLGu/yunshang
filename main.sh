#!/bin/bash

## vars
PROJECT=github.com/itang/yunshang/main
task=$1

## functions
function do_godoc() {
  pid=`pgrep godoc`
  if [ "$pid" != "" ]; then
    echo "kill $pid && godoc -http=:8080 &"
    kill $pid
    echo "> godoc at http://localhost:8080"
    godoc -http=:8080 &
  fi
}

function do_dev_task() {
  do_godoc;
  revel run ${PROJECT} dev
}

## main
(
cd main;
case $task in
    ""|run|dev) do_dev_task;;
          prod) revel run ${PROJECT} prod;;
          test) revel test ${PROJECT} dev;;
       package) revel package ${PROJECT} ;;
           fmt) go fmt ${PROJECT}/...;;
             *) revel $task ${PROJECT} ;;
esac
)
