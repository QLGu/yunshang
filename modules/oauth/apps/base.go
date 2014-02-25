package apps

import (
	"encoding/base64"

	"github.com/itang/yunshang/modules/oauth"
	"github.com/lunny/xorm"
)

type BaseProvider struct {
	App            oauth.Provider
	ClientId       string
	ClientSecret   string
	Scope          string
	AuthURL        string
	TokenURL       string
	RedirectURL    string
	AccessType     string
	ApprovalPrompt string
}

func (p *BaseProvider) getBasicAuth() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(p.ClientId+":"+p.ClientSecret))
}

func (p *BaseProvider) GetConfig() *oauth.Config {
	return &oauth.Config{
		ClientId:       p.ClientId,
		ClientSecret:   p.ClientSecret,
		Scope:          p.Scope,
		AuthURL:        p.AuthURL,
		TokenURL:       p.TokenURL,
		RedirectURL:    p.RedirectURL,
		AccessType:     p.AccessType,
		ApprovalPrompt: p.ApprovalPrompt,
	}
}

func (p *BaseProvider) CanConnect(session *xorm.Session, tok *oauth.Token, userSocial *oauth.UserSocial) (bool, error) {
	identify, err := p.App.GetIndentify(tok) //获取登录用户标识符
	if err != nil {                          // 如果出错， 则表示不能连接
		return false, err
	}
	ok, err := session.Where("identify=? and type=?", identify, p.App.GetType()).Get(userSocial)
	if !ok || err != nil { // 用户连接数据不存在， 表示可以进行连接
		return true, nil
	} else if err == nil {
		return false, nil
	} else {
		return false, err
	}
}
