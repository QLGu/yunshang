package entity

import (
	"time"
)

var (
	ACS_S  = "site"               //站点设置
	ACS_SB = full(ACS_S, "basic") //基本设置
)

var DefaultAppConfs = []AppConfig{
	NewAppConfig("站点名称", full(ACS_SB, "name"), "yunshang", "整个网站的名称"),
	NewAppConfig("默认标题", full(ACS_SB, "title"), "请修改网站名称", "整个站点默认的标题(title)，搜索引擎优化使用"),
	NewAppConfig("默认关键字", full(ACS_SB, "keywords"), "", "页面meta标签里的关键字信息(keywords)，搜索引擎优化使用"),
	NewAppConfig("默认描述信息", full(ACS_SB, "description"), "", "页面meta标签对关键字内容的描述(description)，搜索引擎优化使用"),
	NewAppConfig("ICP/IP/域名备案", full(ACS_SB, "icp"), "公司 xxx ICP备xxxx号", "例如：xxx网络有限公司 粤ICP备00000001号"),
	NewAppConfig("版权信息", full(ACS_SB, "support"), "德福泰-deftype.com", ""),
	NewAppConfig("站点统计", full(ACS_SB, "stat"), "false", "提示：关闭后站点统计将不可用。"),
	NewAppConfig("第三方统计代码", full(ACS_SB, "statcode"), "", ""),
}

func NewAppConfig(name, key, value, description string) AppConfig {
	return AppConfig{Name: name, Key: key, Value: value, Description: description}
}

type AppConfig struct {
	Id int64 `json:"id"`

	Name  string `json:"name"`                       //命名
	Key   string `json:"key"`                        //KEY
	Value string `xorm:"varchar(2000)" json:"value"` //VALUE

	Description string `json:"description"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

func full(s1 string, s2 string) string {
	return s1 + "." + s2
}
