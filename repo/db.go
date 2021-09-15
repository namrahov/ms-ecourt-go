package repo

import (
	"database/sql"
	"fmt"
	"github.com/namrahov/ms-ecourt-go/config"
	log "github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var Conn *sql.DB
var err error

func InitDb() {
	/*Db = pg.Connect(&pg.Options{
		Addr:        config.Props.DbHost + ":" + config.Props.DbPort,
		User:        config.Props.DbUser,
		Password:    config.Props.DbPass,
		Database:    config.Props.DbName,
		PoolSize:    5,
		DialTimeout: 1 * time.Minute,
		MaxRetries:  2,
		MaxConnAge:  15 * time.Minute,
	})*/

	// connect to a database
	Conn, err = sql.Open("pgx",
		fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
			config.Props.DbHost, config.Props.DbPort, config.Props.DbName, config.Props.DbUser, config.Props.DbPass))
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect: %v\n", err))
	}

	log.Println("Connected to database!")

	// test my connection
	err = Conn.Ping()
	if err != nil {
		log.Fatal("Cannot ping database!")
	}

	log.Println("Pinged database!")

	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect %v\n", err))
	}

	log.Println("Connected to db")

	err = Conn.Ping()
	if err != nil {
		log.Fatal("Can not ping database")
	}
	fmt.Println("Pinged database!")

}
