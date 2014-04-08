package migrations

import (
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
