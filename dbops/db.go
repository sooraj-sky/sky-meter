package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "postgress"
	dbname   = "postgres"
)

func main() {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	updatedTime := time.Now()

	insertData := `insert into "endpoint"("end_point_name", "added", "modified") values($1, $2, $3)`
	_, e := db.Exec(insertData, "http://skywalks.in", updatedTime, updatedTime)
	CheckError(e)

	CreateTable := `CREATE TABLE IF NOT EXISTS endpoints (
		id  SERIAL PRIMARY KEY,
		end_point_name TEXT,
		added timestamp default NULL,
		modified timestamp default NULL
		);
		`

	_, t := db.Exec(CreateTable)
	CheckError(t)

	// check db
	err = db.Ping()
	CheckError(err)

    var n int64
    errs := db.QueryRow("select 1 from information_schema.tables where table_name=$1", "endpointk").Scan(&n)
    fmt.Println(errs)



	fmt.Println("Connected!")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
