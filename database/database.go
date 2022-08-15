package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Initializes the db conection
func StartDatabase() {
	connstring := "postgres://root:password@db:5432/emotitron_db?sslmode=disable"

	dbconn, err := sql.Open("postgres", connstring)
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
	query := `CREATE TABLE IF NOT EXISTS customcommands(command_id SERIAL PRIMARY KEY, command_name TEXT, command_result TEXT)`

	createTable(query)
}

// Creates the reminders table
func createRemindersTable() {
	query := `CREATE TABLE IF NOT EXISTS reminders(reminder_id SERIAL PRIMARY KEY, user_id TEXT, future BIGINT, reminder_text TEXT, completed BOOLEAN)`

	createTable(query)
}

// Creates the responses table
func createResponsesTable() {
	query := `CREATE TABLE IF NOT EXISTS responses(response_id SERIAL PRIMARY KEY, response_trigger TEXT, response_value TEXT)`

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
