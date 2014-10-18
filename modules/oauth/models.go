package oauth

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lunny/xorm"
)

const (
	startType SocialType = iota
	SocialGithub
	SocialGoogle
	SocialWeibo
	SocialQQ
	SocialQQT
	SocialDropbox
	SocialFacebook
	endType
)

var types []SocialType

func GetAllTypes() []SocialType {
	if types == nil {
		types = make([]SocialType, int(endType)-1)
		for i, _ := range types {
			types[i] = SocialType(i + 1)
		}
	}
	return types
}

type SocialType int

func (s SocialType) Available() bool {
	if s > startType && s < endType {
		return true
	}
	return false
}

func (s SocialType) Name() string {
	if p, ok := GetProviderByType(s); ok {
		return p.GetName()
	}
	return ""
}

func (s SocialType) NameLower() string {
	return strings.ToLower(s.Name())
}

type SocialTokenField struct {
	*Token
}

func (e *SocialTokenField) String() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func (e *SocialTokenField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case string:
		return json.Unmarshal([]byte(d), e)
	default:
		return fmt.Errorf("<SocialTokenField.SetRaw> unknown value `%v`", value)
	}
	return nil
}

func (e *SocialTokenField) RawValue() string {
	return e.String()
}

type UserSocial struct {
	Id       int64
	Uid      int64 `xorm:"index"`
	Identify string
	Type     SocialType `xorm:"index"`
	RawData  string     `xorm:"varchar(1000)"`

	Session *xorm.Session `xorm:"-"`
}

func (e *UserSocial) Data() (st SocialTokenField) {
	_ = json.Unmarshal([]byte(e.RawData), &st)
	return
}

func (e *UserSocial) Save() (err error) {
	if e.Id == 0 {
		_, err = e.Session.Insert(e)
	} else {
		_, err = e.Session.Id(e.Id).Update(e)
	}
	return
}

func (e *UserSocial) Token() (*Token, error) {
	return e.Data().Token, nil
}

func (e *UserSocial) PutToken(token *Token) error {
	if token == nil {
		return fmt.Errorf("token must be not nil")
	}

	changed := false
	st := e.Data()
	if st.Token == nil {
		st.Token = token
		e.RawData = st.String()
		changed = true
	} else {
		if len(token.AccessToken) > 0 && token.AccessToken != st.AccessToken {
			st.AccessToken = token.AccessToken
			e.RawData = st.String()
			changed = true
		}
		if len(token.RefreshToken) > 0 && token.RefreshToken != st.RefreshToken {
			st.RefreshToken = token.RefreshToken
			e.RawData = st.String()
			changed = true
		}
		if len(token.TokenType) > 0 && token.TokenType != st.TokenType {
			st.TokenType = token.TokenType
			e.RawData = st.String()
			changed = true
		}
		if !token.Expiry.IsZero() && token.Expiry != st.Expiry {
			st.Expiry = token.Expiry
			e.RawData = st.String()
			changed = true
		}
	}

	if changed && e.Id > 0 {
		_, err := e.Session.Id(e.Id).Update(e)
		return err
	}

	return nil
}

func (e *UserSocial) Insert() error {
	if _, err := e.Session.Insert(e); err != nil {
		return err
	}
	return nil
}

func (e *UserSocial) Read(fields ...string) error {
	if err := e.Session.Id(e.Id).Cols(fields...).Find(e); err != nil {
		return err
	}
	return nil
}

func (e *UserSocial) Update(fields ...string) error {
	if _, err := e.Session.Id(e.Id).Cols(fields...).Update(e); err != nil {
		return err
	}
	return nil
}

func (e *UserSocial) Delete() error {
	if _, err := e.Session.Id(e.Id).Delete(e); err != nil {
		return err
	}
	return nil
}
