package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Initializes the db conection
func StartDatabase() {
	dbconn, err := sql.Open("mysql", "root:password@tcp(db:3306)/emotitron_db")
	if err != nil {
		log.Fatal(err)
	}
	db = dbconn
}

// Cleanly closes the db connection
func ShutdownDatabase() {
	db.Close()
}

// Pings the database
func PingDatabase() {
	err := db.Ping()
	if err != nil {
		fmt.Println("Error pinging database", err)
		return
	}

	fmt.Println("Pinged the database")
}

func CreateTable() {
	query := `CREATE TABLE IF NOT EXISTS customcommands(command_id int primary key auto_increment, command_name text, command_result text)`

	// Create a 5 second timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating customcommands table", err)
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return
	}

	log.Printf("Rows affected when creating table: %d", rows)
}
