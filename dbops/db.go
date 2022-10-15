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
var psqlconn string

type errorString struct {
    s string
}

func main() {
	// connection string
	psqlconn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database

	db, err := sql.Open("postgres", psqlconn)
	if err != nil{
		fmt.Println(err)
	}

	// check db
	pingerr := db.Ping()
	if pingerr != nil{
		fmt.Println(err)
	}

	// close database
	defer db.Close()

	updatedTime := time.Now()

	enpoint_table := "endpoint"

	Tableqry := "`"+"CREATE TABLE IF NOT EXISTS"+enpoint_table+"(id  SERIAL PRIMARY KEY,end_point_name TEXT UNIQUE,added timestamp default NULL,modified timestamp default NULL);"+"`"

	CheckTable(db, Tableqry, enpoint_table)

	insertData := `insert into "endpoint"("end_point_name", "added", "modified") values($1, $2, $3)`
	_, e := db.Exec(insertData, "http://skywalks.in", updatedTime, updatedTime)
	CheckError(e, "error adding endoint")


}

func CheckError(err error, errordata string) {

	if err != nil {
	error_time := time.Now()
	db, errs := sql.Open("postgres", psqlconn)
	watchDogQuery := `CREATE TABLE IF NOT EXISTS watchdog (id  SERIAL PRIMARY KEY, operation TEXT, error TEXT, time timestamp default NULL);`
	fmt.Println(errs)
	CheckTable(db, watchDogQuery, "watchdog")
	stringError := err.Error()
	insertError := `insert into "watchdog"("operation","error", "time") values($1, $2, $3)`
	_, e := db.Exec(insertError, errordata, stringError, error_time)
	fmt.Println(e)

	}
	

}

func (e *errorString) Errorstrng() string {
    return e.s
}


func CheckTable(db *sql.DB, Tableqry,Tableqryname string) {
	var n int64 
	errs := db.QueryRow("select 1 from information_schema.tables where table_name=$1", Tableqryname).Scan(&n)
	if errs == nil {

	} else {
		_, t := db.Exec(Tableqry)
		CheckError(t, "check table error")

	}

}
