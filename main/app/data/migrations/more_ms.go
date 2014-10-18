package migrations

import (
	//. "github.com/itang/gotang"
	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/data/migrates"
	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/modules/db"
	"github.com/lunny/xorm"
)

func init() {
	migrates.DataIniter.RegistMigration(m_appConfig())
	migrates.DataIniter.RegistMigration(m_app_params_data())
	migrates.DataIniter.RegistMigration(m_comments_username())
	migrates.DataIniter.RegistMigration(m_product_appConfig())
	migrates.DataIniter.RegistMigration(m_contact_appConfig())
	migrates.DataIniter.RegistMigration(m_links_appConfig())
	migrates.DataIniter.RegistMigration(m_host_appConfig())
	migrates.DataIniter.RegistMigration(m_tags_category())
	migrates.DataIniter.RegistMigration(m_morecontact_appConfig())

	migrates.DataIniter.RegistMigration(m_inner_tags_user())
	migrates.DataIniter.RegistMigration(m_mail_appConfig())
	migrates.DataIniter.RegistMigration(m_alipay_appConfig())
	migrates.DataIniter.RegistMigration(m_bank_appConfig())

	migrates.DataIniter.RegistMigration(m_open_appConfig())

	migrates.DataIniter.RegistMigration(m_shippings())
	migrates.DataIniter.RegistMigration(m_huishouNews())

	migrates.DataIniter.RegistMigration(m_Feedback())

	migrates.DataIniter.RegistMigration(m_shippings_2())

	migrates.DataIniter.RegistMigration(order_payamount())

	migrates.DataIniter.RegistMigration(order_payamount_initdata())
	migrates.DataIniter.RegistMigration(m_qqt_appConfig())
}

func m_appConfig() migrates.Migration {
	return migrates.Migration{
		Name: "app_config",
		Do: func(session *xorm.Session) error {
			db.Engine.Sync(&entity.AppConfig{})
			return models.NewAppConfigService(session).InitData()
		},
	}
}

func m_app_params_data() migrates.Migration {
	return migrates.Migration{
		Name: "m_app_params_data",
		Do: func(session *xorm.Session) error {
			db.Engine.Sync(&entity.AppParams{})
			return nil
		},
	}
}

func m_comments_username() migrates.Migration {
	return migrates.Migration{
		Name: "m_comments_username",
		Do: func(session *xorm.Session) error {
			db.Engine.Sync(&entity.Comment{})
			return nil
		},
	}
}

func m_product_appConfig() migrates.Migration {
	return migrates.Migration{
		Name: "m_product_appConfig",
		Do: func(session *xorm.Session) error {
			appApi := models.NewAppConfigService((session))
			for _, o := range entity.ProductAppConfs {
				appApi.SaveOrUpdateConfigObject(o)
			}
			return nil
		},
	}
}

func m_contact_appConfig() migrates.Migration {
	return migrates.Migration{
		Name: "m_contact_appConfig",
		Do: func(session *xorm.Session) error {
			appApi := models.NewAppConfigService((session))
			for _, o := range entity.ContactAppConfs {
				appApi.SaveOrUpdateConfigObject(o)
			}
			return nil
		},
	}
}

func m_links_appConfig() migrates.Migration {
	return migrates.Migration{
		Name: "m_links_appConfig",
		Do: func(session *xorm.Session) error {
			appApi := models.NewAppConfigService((session))
			for _, o := range entity.LinksAppConfs {
				appApi.SaveOrUpdateConfigObject(o)
			}
			return nil
		},
	}
}

func m_host_appConfig() migrates.Migration {
	return migrates.Migration{
		Name: "m_host_appConfig",
		Do: func(session *xorm.Session) error {
			appApi := models.NewAppConfigService((session))
			for _, o := range entity.HostAppConfs {
				appApi.SaveOrUpdateConfigObject(o)
			}
			return nil
		},
	}
}

func m_tags_category() migrates.Migration {
	return migrates.Migration{
		Name: "m_tags_category",
		Do: func(session *xorm.Session) error {
			db.Engine.Sync(&entity.ProductCategory{})
			return nil
		},
	}
}

func m_morecontact_appConfig() migrates.Migration {
	return migrates.Migration{
		Name: "m_morecontact_appConfig",
		Do: func(session *xorm.Session) error {
			appApi := models.NewAppConfigService((session))
			for _, o := range entity.MoreContactAppConfs {
				appApi.SaveOrUpdateConfigObject(o)
			}
			return nil
		},
	}
}

func m_inner_tags_user() migrates.Migration {
	return migrates.Migration{
		Name: "m_inner_tags_user",
		Do: func(session *xorm.Session) error {
			db.Engine.Sync(&entity.User{})
			models.NewUserService(session).SaveUserRole(1, "#超级管理员")
			return nil
		},
	}
}

func m_mail_appConfig() migrates.Migration {
	return migrates.Migration{
		Name: "m_mail_appConfig",
		Do: func(session *xorm.Session) error {
			appApi := models.NewAppConfigService((session))
			for _, o := range entity.MailAppConfs {
				appApi.SaveOrUpdateConfigObject(o)
			}
			return nil
		},
	}
}

func m_alipay_appConfig() migrates.Migration {
	return migrates.Migration{
		Name: "m_alipay_appConfig",
		Do: func(session *xorm.Session) error {
			appApi := models.NewAppConfigService((session))
			for _, o := range entity.AlipayAppConfs {
				appApi.SaveOrUpdateConfigObject(o)
			}
			return nil
		},
	}
}

func m_bank_appConfig() migrates.Migration {
	return migrates.Migration{
		Name: "m_bank_appConfig",
		Do: func(session *xorm.Session) error {
			db.Engine.Sync(&entity.Bank{})

			ps := []entity.Bank{
				{Name: "ICBCBTB", Description: "中国工商银行(B2B)", Enabled: true},
				{Name: "ABCBTB", Description: "中国农业银行(B2B)", Enabled: true},
				{Name: "CCBBTB", Description: "中国建设银行(B2B)", Enabled: true},
				{Name: "SPDBB2B", Description: "上海浦东发展银行(B2B)", Enabled: true},
				{Name: "BOCBTB", Description: "中国银行(B2B)", Enabled: true},
				{Name: "CMBBTB", Description: "招商银行(B2B)", Enabled: true},
				{Name: "BOCB2C", Description: "中国银行", Enabled: true},
				{Name: "ICBCB2C", Description: "中国工商银行", Enabled: true},
				{Name: "CMB", Description: "招商银行", Enabled: true},
				{Name: "CCB", Description: "中国建设银行", Enabled: true},
				{Name: "ABC", Description: "中国农业银行", Enabled: true},
				{Name: "SPDB", Description: "上海浦东发展银行", Enabled: true},
				{Name: "CIB", Description: "兴业银行", Enabled: true},
				{Name: "GDB", Description: "广发银行", Enabled: true},
				{Name: "CMBC", Description: "中国民生银行", Enabled: true},
				{Name: "CITIC", Description: "中信银行", Enabled: true},
				{Name: "HZCBB2C", Description: "杭州银行", Enabled: true},
				{Name: "CEBBANK", Description: "中国光大银行", Enabled: true},
				{Name: "SHBANK", Description: "上海银行", Enabled: true},
				{Name: "NBBANK", Description: "宁波银行", Enabled: true},
				{Name: "SPABANK", Description: "平安银行", Enabled: true},
				{Name: "BJRCB", Description: "北京农村商业银行", Enabled: true},
				{Name: "FDB", Description: "富滇银行", Enabled: true},
				{Name: "POSTGC", Description: "中国邮政储蓄银行", Enabled: true},
				{Name: "abc1003", Description: "visa", Enabled: true},
				{Name: "abc1004", Description: "master", Enabled: true},
			}
			_, err := session.Insert(ps)
			gotang.AssertNoError(err, "m_bank_appConfig")
			return nil
		},
	}
}

func m_open_appConfig() migrates.Migration {
	return migrates.Migration{
		Name: "m_open_appConfig",
		Do: func(session *xorm.Session) error {
			appApi := models.NewAppConfigService((session))
			for _, o := range entity.OpenLoginAppConfs {
				appApi.SaveOrUpdateConfigObject(o)
			}
			return nil
		},
	}
}

func m_shippings() migrates.Migration {
	return migrates.Migration{
		Name: "m_shippings",
		Do: func(session *xorm.Session) error {
			ps := []entity.Shipping{
				{Name: "优速快递", Description: "", Enabled: true},
				{Name: "全一快递", Description: "", Enabled: true},
				{Name: "EMS快递", Description: "", Enabled: true},
			}
			_, err := session.Insert(ps)
			gotang.AssertNoError(err, "")

			return nil
		},
	}
}

func m_shippings_2() migrates.Migration {
	return migrates.Migration{
		Name: "m_shippings2",
		Do: func(session *xorm.Session) error {
			ps := &entity.Shipping{Name: "宅急送", Description: "", Enabled: true}
		
			_, err := session.Insert(ps)
			gotang.AssertNoError(err, "")

			return nil
		},
	}
}

func m_huishouNews() migrates.Migration {
	return migrates.Migration{
		Name: "m_huishouNews",
		Do: func(session *xorm.Session) error {
			s := models.NewNewsService(session)
			articles := []entity.News{
				{Title: "电子元件回收", CategoryId: 9, Enabled: true},
			}
			for _, e := range articles {
				_, err := s.SaveNews(e)
				gotang.AssertNoError(err, "")
			}
			return nil
		},
	}
}

func m_Feedback() migrates.Migration {
	return migrates.Migration{
		Name: "m_Feedback",
		Do: func(session *xorm.Session) error {
			db.Engine.Sync(&entity.Feedback{})
			return nil
		},
	}
}

func order_payamount() migrates.Migration {
	return migrates.Migration{
		Name: "order_payamount",
		Do: func(session *xorm.Session) error {
			db.Engine.Sync(&entity.Order{})
			return nil
		},
	}
}

func order_payamount_initdata() migrates.Migration {
	return migrates.Migration{
		Name: "order_payamount_initdata",
		Do: func(session *xorm.Session) error {
			db.Engine.Exec("update t_order set pay_amount = amount")
			return nil
		},
	}
}


func m_qqt_appConfig() migrates.Migration {
	return migrates.Migration{
		Name: "m_qqt_appConfig",
		Do: func(session *xorm.Session) error {
			appApi := models.NewAppConfigService((session))
			for _, o := range entity.OpenLoginQQTAppConfs {
				appApi.SaveOrUpdateConfigObject(o)
			}
			return nil
		},
	}
}
