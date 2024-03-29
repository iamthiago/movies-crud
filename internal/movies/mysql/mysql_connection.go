package mysql

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func GetMySQLDB() (db *sql.DB, err error) {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "movies",
		AllowNativePasswords: true,
	}

	db, err = sql.Open("mysql", cfg.FormatDSN())

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	return
}
