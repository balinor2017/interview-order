package repository

import (
	"github.com/interview-order/config"
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

	dbUser = config.MustGetString("mysql.user")
	dbPassword = config.MustGetString("mysql.password")
	dbName = config.MustGetString("mysql.dbname")
	dbHost = config.MustGetString("mysql.host")
	dbPort = config.MustGetString("mysql.port")
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
