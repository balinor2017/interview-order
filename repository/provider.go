package repository

import (
	"github.com/balinor2017/interview-order/config"
	log "github.com/sirupsen/logrus"
)

var (
	psqlDataSource *DataSource
	dbName         string
	dbUser         string
	dbPassword     string
	dbHost         string
	dbPort         string
	redisMain      string
)

func InitDbFactory() error {

	dbUser = config.MustGetString("database.user")
	dbPassword = config.MustGetString("database.password")
	dbName = config.MustGetString("database.dbname")
	dbHost = config.MustGetString("database.host")
	dbPort = config.MustGetString("database.port")
	redisMain = config.MustGetString("cache.redis_main")

	var err error
	psqlDataSource, err = NewDatabaseConnection(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		return err
	}
	log.Info("Init Postgresql Connection Started")
	InitRedis(redisMain)
	log.Info("Init Redis Connection Started")
	return nil
}
