package apps

import (
	"encoding/base64"

	"github.com/lunny/xorm"

	"github.com/itang/yunshang/modules/oauth"
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
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(p.ClientId + ":" + p.ClientSecret))
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

func (p *BaseProvider) CanConnect(session *xorm.Session, tok *oauth.Token, useroauth *oauth.UserSocial) (bool, error) {
	identify, err := p.App.GetIndentify(tok)
	if err != nil {
		return false, err
	}

	if ok, err := session.Where("identify=? and type=?", identify, p.App.GetType()).Get(useroauth); !ok || err != nil {
		return true, nil
	} else if err == nil {
		return false, nil
	} else {
		return false, err
	}
}
