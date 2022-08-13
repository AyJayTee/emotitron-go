package database

import (
	"context"
	"log"
	"time"
)

type CustomCommand struct {
	Name   string
	Result string
}

// Inserts a custom command into the database
func InsertCustomCommand(command CustomCommand) error {
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
	_, err = stmt.ExecContext(ctx, command.Name, command.Result)
	if err != nil {
		log.Printf("Error %s when executing sql query", err)
		return err
	}
	log.Println("Command created:", command.Name)
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

	log.Println("Fetched result for command:", commandName)

	return result, nil
}

// Gets all custom commands in the database
func GetAllCustomCommandNames() ([]CustomCommand, error) {
	query := `SELECT * FROM customcommands`

	// Create a 5 second timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	// Prepare the statement
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing sql query", err)
		return nil, err
	}

	// Scan the result to the variable
	var result []CustomCommand
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		log.Printf("Error %s executing sql query", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var cc CustomCommand
		if err := rows.Scan(&id, &cc.Name, &cc.Result); err != nil {
			log.Printf("Error %s accessing row values", err)
			return nil, err
		}
		result = append(result, cc)
	}

	return result, nil
}

// Removes a custom command from the database
func RemoveCustomCommand(commandName string) error {
	query := `DELETE FROM customcommands WHERE command_name = ?`

	// Create a 5 second timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	// Prepare the statement
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing sql query", err)
		return err
	}
	defer stmt.Close()

	// Execute the query
	_, err = stmt.ExecContext(ctx, commandName)
	if err != nil {
		log.Printf("Error %s when executing the sql query", err)
		return err
	}

	log.Println("Command deleted:", commandName)
	return nil
}
