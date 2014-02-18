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

        CREATE DATABASE yunshangdb OWNER dbuser;
        GRANT ALL PRIVILEGES ON DATABASE yunshangdb to dbuser;

        \c

    psql -U dbuser -d yunshangdb -h 127.0.0.1 -p 5432

        \d
        \d tablename
        \c
        \conninfo

### ref docs

[PostgresSQL][]

[PostgresSQL]: file:///home/itang/resources/postgres/PostgreSQL%E6%96%B0%E6%89%8B%E5%85%A5%E9%97%A8%20-%20%E9%98%AE%E4%B8%80%E5%B3%B0%E7%9A%84%E7%BD%91%E7%BB%9C%E6%97%A5%E5%BF%97.html "PostgresSQL"

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

* 更名栏位 
ALTER TABLE usertbl RENAME COLUMN signupdate TO signup;

* 删除栏位 
ALTER TABLE user_tbl DROP COLUMN email;

* 表格更名 
ALTER TABLE usertbl RENAME TO backuptbl;

* 删除表格 
DROP TABLE IF EXISTS backup_tbl;