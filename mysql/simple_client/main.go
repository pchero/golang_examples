package main

// go run main.go -dsn "testid:testpassword@tcp(127.0.0.1:3306)/mysql"

import (
	"database/sql"
	"flag"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var dsn = flag.String("dsn", "root:root@localhost/test", "target database dsn")
var query = flag.String("query", "select * from user", "query")

func main() {
	flag.Parse()

	db, err := sql.Open("mysql", *dsn)
	if err != nil {
		log.Errorf("Could not access to database. err: %v", err)
		return
	}

	if err := db.Ping(); err != nil {
		log.Errorf("Wrong reponse. err: %v", err)
		return
	}

	// stmt, err := db.Prepare("select User from user;")
	stmt, err := db.Prepare(*query)
	if err != nil {
		log.Errorf("Could not prepare the stmt. err: %v", err)
		return
	}

	row, err := stmt.Query()
	if err != nil {
		log.Errorf("Could not query. err: %v", err)
		return
	}

	var res string
	for row.Next() {
		row.Scan(&res)
		log.Infof("res: %s", res)
	}

	return
}
