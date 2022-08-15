package database

import (
	"context"
	"log"
	"strconv"
	"time"
)

type Reminder struct {
	Id        int
	UserID    string
	Future    int
	Text      string
	Completed bool
}

// Inserts a reminder into the database
func InsertReminder(reminder Reminder) error {
	query := `INSERT INTO reminders (user_id, future, reminder_text, completed) VALUES (`
	query += ("'" + reminder.UserID + "', " + strconv.Itoa(reminder.Future) + ", '" + reminder.Text + "', " + strconv.FormatBool(reminder.Completed))
	query += ")"

	// Create a 5 second timeout
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
	log.Println("Reminder created for user:", reminder.UserID)
	return nil
}

// Removes all reminders belonging to the passed user id
func RemoveReminders(userId string) error {
	// Build the query
	query := `DELETE FROM reminders WHERE user_id = `
	query += "'" + userId + "'"

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

	// Execute the statement
	res, err := stmt.ExecContext(ctx)
	if err != nil {
		log.Printf("Error %s when executing sql query", err)
		return err
	}

	// Get rows affected
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
	}

	log.Printf("Rows affected when deleting reminders: %d", rows)

	return nil
}

// Removes a single reminder from the database
func RemoveRemdinder(id int) error {
	// Build hte query
	query := `DELETE FROM reminders WHERE reminder_id = `
	query += strconv.Itoa(id)

	// Create a 5 second timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	// Prepare the statement
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing sql query", err)
		return err
	}

	// Execute the statement
	res, err := stmt.ExecContext(ctx)
	if err != nil {
		log.Printf("Error %s when executing sql query", err)
		return err
	}

	// Get rows affected
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
	}

	log.Printf("Rows affected when deleting reminders: %d", rows)

	return nil
}

// Gets all expired reminders
func GetExpiredReminders(currentTime int64) ([]Reminder, error) {
	// Build the query
	query := `SELECT reminder_id, user_id, future, reminder_text FROM reminders WHERE future <= `
	query += strconv.Itoa(int(currentTime))

	// Create a 5 second timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	// Prepare the statement
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing sql query", err)
		return nil, err
	}
	defer stmt.Close()

	// Scan the result to a variable
	var result []Reminder
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		log.Printf("Error %s executing sql query", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var reminder Reminder
		if err := rows.Scan(&reminder.Id, &reminder.UserID, &reminder.Future, &reminder.Text); err != nil {
			log.Printf("Error %s accessing row values", err)
			return nil, err
		}
		result = append(result, reminder)
	}

	return result, nil
}
