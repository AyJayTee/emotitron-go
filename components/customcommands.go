package components

import (
	"errors"
)

type customCommand struct {
	name   string
	result string
}

// Adds a custom command to the database
func AddCommand(args []string) (string, error) {
	// Verify that the args are of correct format
	if len(args) != 2 {
		return "", errors.New("correct usage is !add [command name] [command value]")
	}

	// Add the command to the database
	err := InsertCustomCommand(customCommand{args[0], args[1]})
	if err != nil {
		return "", err
	}

	return "Command successfully added!", nil
}

// Removes a custom command from the database
func RemoveCommand(args []string) (string, error) {
	return "Command successfully deleted!", nil
}

// Fetches the value of a custom command from the database
func GetCommand(args []string) (string, error) {
	// Verify that the args are of correct format
	if len(args) != 1 {
		return "", errors.New("wrong number of arguments")
	}

	// Fetch the command from the database
	result, err := GetCustomCommandValue(args[0])
	if err != nil {
		return "", err
	}

	return result, nil
}
