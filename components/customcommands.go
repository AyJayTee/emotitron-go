package components

import (
	"errors"

	"github.com/AyJayTee/emotitron-go/database"
)

// Adds a custom command to the database
func AddCommand(args []string) (string, error) {
	// Verify that the args are of correct format
	if len(args) != 2 {
		return "", errors.New("correct usage is !add [command name] [command value]")
	}

	// Add the command to the database
	err := database.InsertCustomCommand(database.CustomCommand{Name: args[0], Result: args[1]})
	if err != nil {
		return "", err
	}

	return "Command successfully added!", nil
}

// Removes a custom command from the database
func RemoveCommand(args []string) (string, error) {
	// Verify that the args are of correct format
	if len(args) != 1 {
		return "", errors.New("correct usage is !remove [command name]")
	}

	// Remove the command from the database
	err := database.RemoveCustomCommand(args[0])
	if err != nil {
		return "", err
	}

	return "Command successfully deleted!", nil
}

// Fetches the value of a custom command from the database
func GetCommand(args []string) (string, error) {
	// Verify that the args are of correct format
	if len(args) != 1 {
		return "", errors.New("wrong number of arguments")
	}

	// Fetch the command from the database
	result, err := database.GetCustomCommandValue(args[0])
	if err != nil {
		return "", err
	}

	return result, nil
}
