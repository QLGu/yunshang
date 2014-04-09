package migrations

import (
	. "github.com/itang/gotang"
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
			_, err := db.Engine.Exec("ALTER TABLE t_app_config ALTER COLUMN value TYPE varchar(4000)")
			AssertNoError(err, "m_contact_appConfig")

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
