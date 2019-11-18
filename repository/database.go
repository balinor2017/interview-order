package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DataSource struct {
	*sqlx.DB
	maxIdleConns     int
	maxOpenConns     int
	maxConnsLifetime time.Duration
}

func NewDatabaseConnectionWithConnectionPool(host string, port string, user string, password string, database string, maxIdle int, maxOpenConnection int) (*DataSource, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s  sslmode=disable", host, port, user, password, database))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	maxLifeTime := 1 * time.Second

	db.SetMaxIdleConns(maxIdle)
	db.SetMaxOpenConns(maxOpenConnection)
	db.SetConnMaxLifetime(maxLifeTime)
	return &DataSource{db, maxIdle, maxOpenConnection, maxLifeTime}, nil
}

func NewDatabaseConnection(host string, port string, user string, password string, database string) (*DataSource, error) {
	return NewDatabaseConnectionWithConnectionPool(host, port, user, password, database, 10, 10)
}
