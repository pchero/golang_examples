package main

// go run main.go -dsn "testid:testpassword@tcp(127.0.0.1:3306)/mysql"

// mysql> desc test;
// +-------+---------+------+-----+---------+-------+
// | Field | Type    | Null | Key | Default | Extra |
// +-------+---------+------+-----+---------+-------+
// | id    | int(11) | NO   | PRI | NULL    |       |
// | data  | json    | YES  |     | NULL    |       |
// +-------+---------+------+-----+---------+-------+
// 2 rows in set (0.02 sec)

// example for select json type data from the table and unmarshaling into map[string]string.

import (
	"database/sql"
	"encoding/json"
	"flag"

	_ "github.com/go-sql-driver/mysql"

	log "github.com/sirupsen/logrus"
)

var dsn = flag.String("dsn", "root:root@localhost/test", "target database dsn")

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

	// insert
	stmtIn, err := db.Prepare("insert into test values(?, ?)")
	if err != nil {
		log.Errorf("Could not prepare. err: %v", err)
		return
	}

	mapData := map[string]string{
		"key1": "value1",
	}
	tmpData, err := json.Marshal(mapData)
	if err != nil {
		log.Errorf("Could not marshal. err: %v", err)
		return
	}

	ret, err := stmtIn.Exec(11, tmpData)
	if err != nil {
		log.Errorf("Could not insert. err: %v", err)
		return
	}
	log.Infof("ret: %v", ret)

	// insert empty map
	stmtIn, err = db.Prepare("insert into test values(?, ?)")
	if err != nil {
		log.Errorf("Could not prepare. err: %v", err)
		return
	}

	var mapData2 map[string]string
	tmpData, err = json.Marshal(mapData2)
	if err != nil {
		log.Errorf("Could not marshal. err: %v", err)
		return
	}

	ret, err = stmtIn.Exec(12, tmpData)
	if err != nil {
		log.Errorf("Could not insert. err: %v", err)
		return
	}
	log.Infof("ret: %v", ret)

	// select
	stmt, err := db.Prepare("select * from test;")
	if err != nil {
		log.Errorf("Could not prepare the stmt. err: %v", err)
		return
	}

	row, err := stmt.Query()
	if err != nil {
		log.Errorf("Could not query. err: %v", err)
		return
	}

	var id int
	var data string
	for row.Next() {
		row.Scan(
			&id,
			&data,
		)
		log.Infof("Scan complete. id: %d, data: %s", id, data)
	}

	var mapString map[string]string
	if err := json.Unmarshal([]byte(data), &mapString); err != nil {
		log.Errorf("Could not unmarshal. err: %v", err)
		return
	}
	log.Infof("map: %s", mapString)

	return
}
