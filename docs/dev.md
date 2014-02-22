DEV
===

## code style

* [Google Code Style](https://code.google.com/p/go-wiki/wiki/Style "google")

## convert a []T to an []interface{}?
* https://code.google.com/p/go-wiki/wiki/InterfaceSlice
* http://golang.org/doc/faq#convert_slice_of_interface

## clear data

        drop table t_company;
        drop table t_company_detail_biz;
        drop table t_company_main_biz;
        drop table t_company_type;
        drop table t_location;
        drop table t_user;
        drop table t_user_detail;
        drop table t_user_level;
        drop table t_user_work_kind;

## go get

go get -u -v github.com/robfig/revel/...
go get -u -v github.com/lib/pg

go get -u -v github.com/nu7hatch/gouuid
go get -u -v github.com/itang/gotang
go get -u -v github.com/lunny/xorm

## question

### revel

* Failed to generate name for field. https://github.com/robfig/revel/issues/343