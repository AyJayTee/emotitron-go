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

	// Initialize all the requried tables
	createCustomCommandsTable()
	createRemindersTable()
	createResponsesTable()
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

// Creates the customcommands table
func createCustomCommandsTable() {
	query := `CREATE TABLE IF NOT EXISTS customcommands(command_id int primary key auto_increment, command_name text, command_result text)`

	createTable(query)
}

// Creates the reminders table
func createRemindersTable() {
	query := `CREATE TABLE IF NOT EXISTS reminders(reminder_id int primary key auto_increment, user_id text, future int, reminder_text text, completed boolean)`

	createTable(query)
}

// Creates the responses table
func createResponsesTable() {
	query := `CREATE TABLE IF NOT EXISTS responses(response_id int primary key auto_increment, response_trigger text, response_value text)`

	createTable(query)
}

// Creates a table from a query
func createTable(query string) error {
	// Create 5 second timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	// Execute the statement
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating table", err)
		return err
	}

	// Print out rows affected
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}

	log.Printf("Rows affected when creating table: %d", rows)

	return nil
}
