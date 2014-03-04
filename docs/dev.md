DEV
===

## Components

1. <http://fancyapps.com/fancybox/>

## code style

* [Google Code Style For Go](https://code.google.com/p/go-wiki/wiki/Style "google")

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

$ ./run goupdate

go get -u -v github.com/itang/yunshang/...
go get -u -v github.com/robfig/revel/revel
go get -u -v github.com/robfig/revel/...
go get -u -v github.com/itang/reveltang/...

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


### online

经核实您的域名：xxxx.com对应的备案服务商非阿里云/万网，因工信部要求域名指向的服务器提供商需与域名备案接入商保持一致，若您的域名备案接入商非阿里云/万网，您在使用阿里云服务器时，会出现网站无法访问的情况。需要将顶级域名备案信息进行接入操作。接入申请提交至初审通过6个小时左右您的网站可以使用。申请接入备案操作指南
<http://help.aliyun.com/guide?spm=0.0.0.0.mv9Wad&helpId=877>


### 管理端界面

matrix-admin03 适合做文档的， FAQ的


### bugs

xorm: [x]

self.session.Id(nil) , Id 会传递

UPDATE "t_user" SET "last_sign_at" = $1, "updated_at" = $2, "_version" = "_version" + 1 WHERE ((id=$3) AND (id=$4)) AND "_version" = $5
[2014-02-25 12:50:53.345810596 +0800 CST 2014-02-25 12:50:53.345837231 +0800 CST 1 3 2]

google email:

<http://www.serversmtp.com/en/limits-of-gmail-smtp-server>
