package repo

import (
	"database/sql"
	"github.com/go-pg/pg"
	_ "github.com/lib/pq"
	"github.com/namrahov/ms-ecourt-go/config"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
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

func MigrateDb() error {
	log.Info("MigrateDb.start")

	connStr := "dbname=" + config.Props.DbName + " user=" + config.Props.DbUser + " password=" + config.Props.DbPass + " host=" + config.Props.DbHost + " port=" + config.Props.DbPort + "  sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	log.Info("Applied ", n, " migrations")
	log.Info("MigrateDb.end")
	return nil
}
