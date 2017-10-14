package db

import (
	"database/sql"
	"log"
	"wxwebupload/tool"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	conf, err := tool.Conf("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	db, err = sql.Open("mysql", conf["db"])
	if err != nil {
		log.Fatal(err)
	}
}
