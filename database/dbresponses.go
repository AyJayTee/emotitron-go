package database

import (
	"context"
	"log"
	"time"
)

type Response struct {
	Trigger  string
	Response string
}

// Inserts a response to the database
func InsertResponse(response Response) error {
	// Build the query
	query := `INSERT INTO responses (response_trigger, response_value) VALUES (`
	query += "'" + response.Trigger + "', '" + response.Response + "')"

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
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		log.Printf("Error %s when executing sql query", err)
		return err
	}
	log.Println("Response created:", response.Trigger)
	return nil
}

// Removes a response from the database
func RemoveResponse(responseTrigger string) error {
	// Build the query
	query := `DELETE FROM responses WHERE response_trigger = `
	query += "'" + responseTrigger + "'"

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
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		log.Printf("Error %s when executing the sql query", err)
		return err
	}

	log.Println("Response deleted:", responseTrigger)
	return nil
}

// Updates the trigger of a repsonse
func UpdateResponseTrigger(oldTrigger string, newTrigger string) error {
	// Build the query
	query := `UPDATE responses SET response_trigger = `
	query += "'" + newTrigger + "'"
	query += " WHERE response_trigger = "
	query += "'" + oldTrigger + "'"

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
	res, err := stmt.ExecContext(ctx)
	if err != nil {
		log.Printf("Error %s when executing the sql query", err)
		return err
	}

	// Get the rows affected
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}

	log.Printf("Rows affected when modifying response: %d", rows)

	return nil
}

// Update the response of a response
func UpdateResponseResponse(trigger string, newRepsonse string) error {
	// Build the query
	query := `UPDATE responses SET response_value = `
	query += "'" + newRepsonse + "'"
	query += " WHERE response_trigger = "
	query += "'" + trigger + "'"

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
	res, err := stmt.ExecContext(ctx)
	if err != nil {
		log.Printf("Error %s when executing the sql query", err)
		return err
	}

	// Get the rows affected
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}

	log.Printf("Rows affected when modifying response: %d", rows)

	return nil
}

// Gets a reponse object based on a trigger
func GetResponse(trigger string) (Response, error) {
	// Build the query
	query := `SELECT response_trigger, response_value FROM responses WHERE response_trigger = `
	query += "'" + trigger + "'"

	// Create a 5 second timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	// Prepare the statment
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing sql query", err)
		return Response{}, err
	}
	defer stmt.Close()

	// Scan the result to a variable
	var result Response
	row := stmt.QueryRowContext(ctx)
	if err := row.Scan(&result.Trigger, &result.Response); err != nil {
		return Response{}, err
	}

	log.Println("Fetched result for response:", trigger)

	return result, nil
}

// Get all stored responses
func GetAllResponses() ([]Response, error) {
	query := `SELECT response_trigger, response_value FROM responses`

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
	var result []Response
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		log.Printf("Error %s executing sql query", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Response
		if err := rows.Scan(&r.Trigger, &r.Response); err != nil {
			log.Printf("Error %s accessing row values", err)
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil
}
