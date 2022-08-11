package components

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func StartDatabase() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/emotitron_db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func ShutdownDatabase(db *sql.DB) {
	db.Close()
}

func PingDatabase(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		fmt.Println("Error pinging database", err)
	}
}
