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
	SocialDropbox
	SocialFacebook
	endType
)

var types []SocialType

func GetAllTypes() []SocialType {
	if types == nil {
		types = make([]SocialType, int(endType)-1)
		for i, _ := range types {
			types[i] = SocialType(i+1)
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

/*
func (e *SocialTokenField) FieldType() int {
	return orm.TypeTextField
}*/

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
	Id       int
	Uid      int              `xorm:"index"`
	Identify string // `orm:"size(200)"`
	Type     SocialType       `xorm:"index"`
	RawData  string `xorm:"varchar(1000)"`

	Session *xorm.Session `xorm:"-"`
}

func (e *UserSocial) Data() (ret SocialTokenField ){
	 _ = json.Unmarshal([]byte(e.RawData), &ret)
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

	if e.Data().Token == nil {
		d := e.Data()
		d.Token = token
		e.RawData = d.String()
		changed = true
	} else {
		d :=e.Data()
		if len(token.AccessToken) > 0 && token.AccessToken != d.AccessToken {
			d.AccessToken = token.AccessToken
			e.RawData = d.String()
			changed = true
		}
		if len(token.RefreshToken) > 0 && token.RefreshToken != d.RefreshToken {
			d.RefreshToken = token.RefreshToken
			e.RawData = d.String()
			changed = true
		}
		if len(token.TokenType) > 0 && token.TokenType != d.TokenType {
			d.TokenType = token.TokenType
			e.RawData = d.String()
			changed = true
		}
		if !token.Expiry.IsZero() && token.Expiry != d.Expiry {
			d.Expiry = token.Expiry
			e.RawData = d.String()
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

/*
func UserSocials() *xorm.Session {
	return orm.NewOrm().QueryTable("user_social")
}
*/
/*
// Get UserSocials by uid
func GetSocialsByUid(session *xorm.Session, uid int, socialTypes ...SocialType) ([]*UserSocial, error) {
	var userSocials []*UserSocial
	_, err := e.Session.Where("Uid =? and Type in ?", uid).Filter("Type__in", socialTypes).All(&userSocials)
	if err != nil {
		return nil, err
	}
	return userSocials, nil
}
*/
func init() {
	//orm.RegisterModel(new(UserSocial))
}
