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
