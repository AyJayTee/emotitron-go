package database

import (
	"context"
	"log"
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
	query := `INSERT INTO reminders (user_id, future, reminder_text, completed) VALUES (?, ?, ?, ?)`

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
	_, err = stmt.ExecContext(ctx, reminder.UserID, reminder.Future, reminder.Text, reminder.Completed)
	if err != nil {
		log.Printf("Error %s when executing sql query", err)
		return err
	}
	log.Println("Reminder created for user:", reminder.UserID)
	return nil
}

// Removes all reminders belonging to the passed user id
func RemoveReminders(userId string) error {
	query := `DELETE FROM reminders WHERE user_id = ?`

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
	res, err := stmt.ExecContext(ctx, userId)
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
	query := `DELETE FROM reminders WHERE reminder_id = ?`

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
	res, err := stmt.ExecContext(ctx, id)
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
	query := `SELECT reminder_id, user_id, future, reminder_text FROM reminders WHERE future <= ?`

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
	rows, err := stmt.QueryContext(ctx, currentTime)
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
