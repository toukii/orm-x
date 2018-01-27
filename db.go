package orm

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	User     string
	Passwd   string
	Host     string
	Port     int
	Database string
}

type DBStore struct {
	*sql.DB
}

var (
	_db *DBStore
)

func Mysql() *DBStore {
	return _db
}

func SetUpMysql(cfg *Config) error {
	var dsn string

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&autocommit=true&parseTime=True",
		cfg.User, cfg.Passwd, cfg.Host, cfg.Port, cfg.Database)
	var err error
	rawdb, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	_db = &DBStore{
		rawdb,
	}
	return Mysql().Ping()
}
