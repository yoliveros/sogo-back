package handler

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	db, _ = sql.Open("mysql", "root:root@tcp(localhost:3306)/sogo")
}
