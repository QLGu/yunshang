package oauth

import (
	"github.com/robfig/revel"
	"github.com/lunny/xorm"
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
	IsUserLogin(*revel.Controller) (int, bool)
	LoginUser(*revel.Controller, int) (string, error)
}
