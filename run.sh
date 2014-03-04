#!/bin/bash

## vars
PROJECT=github.com/itang/yunshang/main
PROJECT_MODULES=github.com/itang/yunshang/modules
tasks=$*

## functions
function do_godoc() {
    pid=`pgrep godoc`
    if [ "$pid" != "" ]; then
        echo "kill $pid && godoc -http=:8080 &"
        kill $pid
    fi
    echo "> godoc at http://localhost:8080"
    godoc -http=:8080 &
}

function do_dev_task() {
    go version;
    echo "go fmt ${PROJECT}/..." && go fmt ${PROJECT}/...
    echo "go fmt ${PROJECT_MODULES}/..." && go fmt ${PROJECT_MODULES}/...
    do_godoc;
    revel run ${PROJECT} dev
}

function do_dev_sync() {
    ssh root@godocking.com '(cd yunshang;git pull)'
}

function do_push() {
    git add . -A
    git commit -a -m "update"
    git push origin master
}

function do_goupdate() {
    go get -u -v github.com/revel/cmd/revel
    go install github.com/revel/cmd/revel
    go get -u -v github.com/revel/revel/...
    go get -u -v github.com/itang/reveltang/...

    go get -u -v github.com/lib/pq
    go get -u -v github.com/nu7hatch/gouuid
    go get -u -v github.com/itang/gotang
    go get -u -v github.com/lunny/xorm
    go get -u -v github.com/astaxie/beego/httplib
    go get -u -v github.com/go-sql-driver/mysql
    go get -u -v github.com/ungerik/go-mail
}

#####################################################################
## main
(cd main;
if [[ "$tasks" = "" ]]; then
    tasks="dev"
fi

for task in $tasks; do
    case $task in
        "" | run | dev) do_dev_task ;;
        prod) revel run ${PROJECT} prod ;;
        test) revel test ${PROJECT} dev ;;
        package) revel package ${PROJECT} ;;
        fmt) go fmt ${PROJECT}/... ;;
        initdb) go run ../tools/initdb.go ;;
        psql) psql -U dbuser -d yunshangdb -h 127.0.0.1 -p 5432 ;;
        dev-sync) do_dev_sync ;;
        push) do_push ;;
        goupdate) do_goupdate ;;
        *) revel $task ${PROJECT} ;;
    esac
done)
