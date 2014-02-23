DEV
===

## code style

* [Google Code Style](https://code.google.com/p/go-wiki/wiki/Style "google")

## convert a []T to an []interface{}?
* https://code.google.com/p/go-wiki/wiki/InterfaceSlice
* http://golang.org/doc/faq#convert_slice_of_interface

## dev references

* OAuth2: http://blog.yorkxin.org/posts/2013/09/30/oauth2-1-introduction
* zhifubo: https://github.com/yaofangou/open_taobao
* social-auth: https://github.com/beego/social-auth

## clear data

$ ./run initdb

## nginx

install:

    $ sudo apt-get install nginx
    $ ls /etc/nginx

start:

    $ sudo nginx

control:

    $ nginx -s signal # stop quit reload reopen

config:

    $ cat > /etc/nginx/sites-enabled/godocking.conf

    server {
        listen 80;
        server_name  yunshang.godocking.com;
        root         /root/gopath/src/github.com/itang/yunshang/main/public;
        location / {
            proxy_pass http://localhost:9000;
        }
    }

logs:

    $ cat /var/log/nginx/error.log


misc:

haoshuju.wicp.net

## deploy

### postgres

$ sudo apt-get install postgresql-client

$ apt-get install postgresql postgresql-common postgresql-9.1 postgresql-contrib-9.1 # sudo apt-get install postgresql

$ sudo adduser dbuser

$ sudo su - postgres

$ psql

  \password postgres
  IC1

  CREATE USER dbuser WITH PASSWORD 'dbuser';

  CREATE DATABASE yunshangdb OWNER dbuser;

  GRANT ALL PRIVILEGES ON DATABASE yunshangdb to dbuser;

$ psql -U dbuser -d yunshangdb -h 127.0.0.1 -p 5432

### git

apt-get install git

### go

$ mkdir downloads
$ cd downloads
$ wget https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz
$ tar zxvf go


## go get

go get -u -v github.com/itang/yunshang/...
go get -u -v github.com/robfig/revel/revel
go get -u -v github.com/robfig/revel/...

go get -u -v github.com/lib/pg
go get -u -v github.com/nu7hatch/gouuid
go get -u -v github.com/itang/gotang
go get -u -v github.com/lunny/xorm

## question

### revel

* Failed to generate name for field. https://github.com/robfig/revel/issues/343

### account

google email:

    yunshang2014
    re***24

weibo:

http://open.weibo.com/apps

qq:

http://connect.qq.com/manage/index
