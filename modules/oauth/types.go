package oauth

import (
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

// Interface of social Privider
type Provider interface {
	GetConfig() *Config
	GetType() SocialType
	GetName() string
	GetPath() string
	GetIndentify(*Token) (string, error)
	CanConnect(*xorm.Session, *Token, *UserSocial) (bool, error)
}

// Interface of social utils
type SocialAuther interface {
	IsUserLogin(*revel.Controller) (int64, bool)
	LoginUser(*revel.Controller, int64, SocialType) (string, error)
}
