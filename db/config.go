package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	DB, _ = sql.Open("mysql", "root:root@tcp(localhost:3306)/sogo")
}

func DeinitDB() {
	DB.Close()
}
