## 登录云主机

$ ssh root@114.215.189.226

(按提示输入密码登录)

## 准备

    $ apt-get update
    $ apt-get upgrade

    $ echo 'export LC_CTYPE=en_US.UTF-8' >> .profile
    $ echo 'export LC_ALL=en_US.UTF-8' >> .profile
    $ source .profile

    $ mkdir dev-env
    $ mkdir tmp
    $ mkdir gopath

## git

### 安装

    $ apt-get install git

    $ git --version

## go

### 安装

    $ cd tmp
    $ wget https://storage.googleapis.com/golang/go1.2.2.linux-amd64.tar.gz
    $ tar zxvf go1.2.2.linux-amd64.tar.gz -C ~/dev-env/
    $ cd

### 设置

    $ echo 'export GOROOT=~/dev-env/go' >> .profile
    $ echo 'export GOPATH=~/gopath' >> .profile
    $ echo 'export PATH=$GOROOT/bin:$GOPATH/bin:$HOME/bin:$PATH' >> .profile
    $ source .profile

    $ go version

    go version go1.2.2 linux/amd64

## 数据库（postgresql）

### 安装

$ apt-get install postgresql-client

$ apt-get install postgresql postgresql-common postgresql-9.1 postgresql-contrib-9.1

### 设置

#### 数据库设置

1. 限制本机才能连接

    修改/etc/postgresql/9.1/main/postgresql.conf
    将行：
    #listen_addresses = 'localhost'
    去掉'#'
    listen_addresses = 'localhost'

2. restart 数据库

#### 数据库用户

1. 创建用户

$ adduser dbuser

(按提示设置密码, 比如输入dbuser)

2. 切换到postgres用户

$ su - postgres

postgres@xxx:~$

3. 使用psql命令登录PostgreSQL控制台

$ psql

4. 为postgres用户设置一个密码

    \password postgres

    (按提示输入密码)

5. 创建数据库用户dbuser（刚才创建的是Linux系统用户），并设置密码

    CREATE USER dbuser WITH PASSWORD 'dbuser';

    这里指定密码为dbuser

6. 创建用户数据库yunshangdb，并指定所有者为dbuser

    CREATE DATABASE yunshangdb OWNER dbuser ENCODING 'UTF8';

7. yunshangdb数据库的所有权限都赋予dbuser

    GRANT ALL PRIVILEGES ON DATABASE yunshangdb to dbuser;

8. dbuser登录数据库

    psql -U dbuser -d yunshangdb -h 127.0.0.1 -p 5432

    (按提示输入密码,如dbuser)

### 管理

1. 重启

/etc/init.d/postgresql restart


### 参考资料

1. <http://www.pixelite.co.nz/article/installing-and-configuring-postgresql-91-ubuntu-1204-local-drupal-development>
2. <http://www.ruanyifeng.com/blog/2013/12/getting_started_with_postgresql.html>
3. <http://install-things.com/2012/06/06/how-to-install-postgres-9-1-on-ubuntu-12-04-linux/>
4. <http://askubuntu.com/questions/157850/im-trying-to-install-postgresql-on-12-04-and-its-just-not-working>

## nginx

### 安装

$ apt-get install nginx

$ nginx -v

### 设置

    $ cat > /etc/nginx/sites-enabled/keeptops.net.conf

    输入以下内容：

    server {
        listen 80;
        server_name  www.keeptops.net;
        root   workspace/gopath/src/github.com/itang/yunshang/main/public;
        location / {
            proxy_pass http://localhost:9000;
            proxy_set_header            X-real-ip $remote_addr;
            proxy_connect_timeout 120s;
            client_max_body_size 50M;
            client_body_buffer_size 128k;
        }
    }

重启nginx

    $ /etc/init.d/nginx restart

### 管理

    $ /etc/init.d/nginx restart

## yunshang程序

### 安装

$ go get -u -v github.com/itang/yunshang/...

$ ln -s gopath/src/github.com/itang/yunshang yunshang


### 设置


### 管理

    $ cd yunshang

#### 运行：

    $ ./run

#### 停止：

找到对应的pid:

    $ ps -ef | grep yunshang/main

kill进程:

    $ kill x1 x2

#### 重启

按以上步骤，先停止再运行

## 域名

添加A记录， 将www.keeptops.net 域名解析指向到114.215.189.226

那么 通过http://www.keeptops.net 就可以访问网站了

## 通过系统管理员管理网站

  admin（默认密码: computer)用户（超级管理员)登录， 通过“我的管理”，进行登录密码的修改， 用户角色的分配

  角色包括： 管理员 销售 超级管理员 , 分别看到不同的操作菜单， 管理如产品、供应商、文章、系统维护等...

  用户通过自主注册使用此网站。
