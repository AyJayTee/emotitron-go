package components

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

// Inserts a custom command into the database
func InsertCustomCommand(command customCommand) error {
	query := `INSERT INTO customcommands (command_name, command_result) VALUES (?, ?)`

	// Create 5 second timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	// Prepare the statement
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when prepararing sql query", err)
		return err
	}
	defer stmt.Close()

	// Execute the statement
	res, err := stmt.ExecContext(ctx, command.name, command.result)
	if err != nil {
		log.Printf("Error %s when executing sql query", err)
		return err
	}

	// Print the rows affected
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}
	log.Printf("%d command created ", rows)
	return nil
}

// Get the value of a custom command
func GetCustomCommandValue(commandName string) (string, error) {
	query := `SELECT command_result FROM customcommands WHERE command_name = ?`

	// Create a 5 second timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	// Prepare the statment
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing sql query", err)
		return "", err
	}
	defer stmt.Close()

	// Scan the result to a variable
	var result string
	row := stmt.QueryRowContext(ctx, commandName)
	if err := row.Scan(&result); err != nil {
		return "", err
	}

	log.Println("Fetched command for command:", commandName)

	return result, nil
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
