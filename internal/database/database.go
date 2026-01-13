package database

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"todogin/internal/config"
)

type Database struct {
	Conn *sql.DB
}

func DatabaseInit(conf *config.Config) (*Database, error) {
	mconf := mysql.Config{
		Addr                : conf.DBAddr,
		User                : conf.DBUser,
		Passwd              : conf.DBPasswd,
		DBName              : conf.DBName,
		AllowNativePasswords: true,
	}

	conn, err := sql.Open("mysql", mconf.FormatDSN())
	if err != nil {
		return nil, &DatabaseError{"sql.Open fail", err}
	}

	err = conn.Ping()
	if err != nil {
		return nil, &DatabaseError{"Ping fail", err}
	}

	database := Database{Conn: conn} 
	return &database, nil
}
