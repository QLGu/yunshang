YUNSHANG DB
===========

## MySQL

## Create DB
## PostgreSQL

### create DB
 
    $ sudo adduser dbuser

    $ sudo su - postgres
    $ psql

        CREATE USER dbuser WITH PASSWORD 'password';

        CREATE DATABASE yunshangdb OWNER dbuser WITH ENCODING 'UTF8';
        GRANT ALL PRIVILEGES ON DATABASE yunshangdb to dbuser;

        \c

    $ psql -U dbuser -d yunshangdb -h 127.0.0.1 -p 5432

        \d
        \d tablename
        \c
        \conninfo

### ref docs

[PostgresSQL][]

[PostgresSQL]: <http://www.ruanyifeng.com/blog/2013/12/getting_started_with_postgresql.html>
基本的数据库操作，就是使用一般的SQL语言。

* 创建新表 
CREATE TABLE usertbl(name VARCHAR(20), signupdate DATE);

* 插入数据 
INSERT INTO usertbl(name, signupdate) VALUES('张三', '2013-12-22');

* 选择记录 
SELECT * FROM user_tbl;

* 更新数据 
UPDATE user_tbl set name = '李四' WHERE name = '张三';

* 删除记录 
DELETE FROM user_tbl WHERE name = '李四' ;

* 添加栏位 
ALTER TABLE user_tbl ADD email VARCHAR(40);

* 更新结构 
ALTER TABLE usertbl ALTER COLUMN signupdate SET NOT NULL;

* 更改类型
ALTER TABLE t_app_config ALTER COLUMN value TYPE varchar(4000);

* 更名栏位 
ALTER TABLE usertbl RENAME COLUMN signupdate TO signup;

* 删除栏位 
ALTER TABLE user_tbl DROP COLUMN email;

* 表格更名 
ALTER TABLE usertbl RENAME TO backuptbl;

* 删除表格 
DROP TABLE IF EXISTS backup_tbl;