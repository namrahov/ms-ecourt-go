package repo

import (
	"github.com/ms-ecourt-go/config"
	"time"
)

var Db *pg.DB

func InitDb() {
	Db = pg.Connect(&pg.Options{
		Addr:        config.Props.DbHost + ":" + config.Props.DbPort,
		User:        config.Props.DbUser,
		Password:    config.Props.DbPass,
		Database:    config.Props.DbName,
		PoolSize:    5,
		DialTimeout: 1 * time.Minute,
		MaxRetries:  2,
		MaxConnAge:  15 * time.Minute,
	})
}
