package entity

import "time"

type User struct {
	Id          int64
	LoginName   string `xorm:"unique"`
	RealName    string
	NickName    string
	Age         int       `xorm:"index"`
	CreatedAt   time.Time `xorm:created`
	UpdatedAt   time.Time `xorm:updated`
	DataVersion int       `xorm:version '_version'`
	Title       string
	Address     string
	Genre       string
	Area        string
}

/*

drop table t_user;

*/
