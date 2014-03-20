#!/bin/bash

## vars
PROJECT=github.com/itang/yunshang
MAIN=${PROJECT}/main
tasks=$*

## functions

function do_fmt() {
    echo "go fmt ${PROJECT}/..."
    go fmt ${PROJECT}/...
}

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

    do_fmt;

    revel run ${MAIN} dev
}

function do_dev_sync() {
    ssh root@godocking.com '(cd yunshang;git pull)'
}

function do_push() {
    do_fmt;
    git add --all .
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
    go get -u -v github.com/deckarep/golang-set
    go get -u -v github.com/disintegration/imaging
    go get -u -v github.com/chuckpreslar/emission
}

function do_catjs() {
    LIBS="/${PWD}/public/libs"
    MEDIA="/${PWD}/public/media"
    COREJS=$LIBS/core.js
    EXTRAJS=$LIBS/extra.js

    cat $LIBS/stringjs/string.min.js > $COREJS
    cat $LIBS/lodash/lodash.compat.min.js >> $COREJS
    cat $LIBS/moment/moment.min.js >> $COREJS
    cat $LIBS/Ractive/Ractive-legacy.min.js >> $COREJS
    cat $MEDIA/js/jquery-1.10.1.min.js >> $COREJS
    cat $MEDIA/js/jquery-migrate-1.2.1.min.js >> $COREJS

    cat $MEDIA/js/jquery-ui-1.10.1.custom.min.js > $EXTRAJS
    cat $MEDIA/js/bootstrap.min.js >> $EXTRAJS
    cat $MEDIA/js/jquery.slimscroll.min.js >> $EXTRAJS
    cat $MEDIA/js/jquery.blockui.min.js >> $EXTRAJS
    cat $MEDIA/js/jquery.cookie.min.js >> $EXTRAJS
    cat $MEDIA/js/jquery.uniform.min.js >> $EXTRAJS
    cat $MEDIA/js/select2.min.js >> $EXTRAJS
    #wget http://ivaynberg.github.io/select2/select2-3.4.5/select2.js
    #cat select2.js >> $EXTRAJS
    #rm select2.js
    cat $MEDIA/js/jquery.dataTables.js >> $EXTRAJS
    cat $MEDIA/js/ZeroClipboard.js >> $EXTRAJS
    cat $MEDIA/js/TableTools.js >> $EXTRAJS
    cat $MEDIA/js/dataTables.bootstrap.js >> $EXTRAJS
    cat $LIBS/fancybox/lib/jquery.mousewheel-3.0.6.pack.js >> $EXTRAJS
    cat $LIBS/fancybox/source/jquery.fancybox.pack.js >> $EXTRAJS
    cat $MEDIA/js/jquery.form.min.js >> $EXTRAJS
    cat $MEDIA/js/jquery.MultiFile.js >> $EXTRAJS
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
        doc) do_godoc ;;
        prod) revel run ${MAIN} prod ;;
        test) revel test ${MAIN} dev ;;
        package) revel package ${MAIN} ;;
        fmt) do_fmt ;;
        initdb) go run ../tools/initdb.go ;;
        psql) psql -U dbuser -d yunshangdb -h 127.0.0.1 -p 5432 ;;
        dev-sync | deploy) do_dev_sync ;;
        push) do_push ;;
        goupdate) do_goupdate ;;
        catjs) do_catjs ;;
        *) revel $task ${MAIN} ;;
    esac
done)
