package components

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/AyJayTee/emotitron-go/database"
	"github.com/bwmarrin/discordgo"
)

// Adds a custom command to the database
func AddCustomCommand(currentCommands []string, m *discordgo.MessageCreate) (string, error) {
	hasAttachment := false
	args := strings.Split(m.Content[5:], " ")

	// Check for attachments
	if len(m.Attachments) != 0 {
		hasAttachment = true
	}

	// Verify that the args are of correct format
	if hasAttachment {
		if len(args) != 1 {
			return "", errors.New("correct usage is !add [command name] [command value], or !add [command name] with an attachment")
		}
	} else {
		if len(args) != 2 {
			return "", errors.New("correct usage is !add [command name] [command value], or !add [command name] with an attachment")
		}
	}

	// Stop user from adding commands named the same as proper commands
	for _, c := range currentCommands {
		if c == args[0] {
			return "", errors.New("cannot add a command that already exists")
		}
	}

	// Check that command does not already exist
	_, err := database.GetCustomCommandValue(args[0])
	if err == nil {
		return "", errors.New("cannot add a command that already exists")
	}

	// Add the command to the database
	if hasAttachment {
		err := database.InsertCustomCommand(database.CustomCommand{Name: args[0], Result: m.Attachments[0].URL})
		if err != nil {
			return "", err
		}
	} else {
		err := database.InsertCustomCommand(database.CustomCommand{Name: args[0], Result: args[1]})
		if err != nil {
			return "", err
		}
	}

	return "Command successfully added!", nil
}

// Removes a custom command from the database
func RemoveCustomCommand(m *discordgo.MessageCreate) (string, error) {
	args := strings.Split(m.Content[8:], " ")

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
func GetCustomCommand(m *discordgo.MessageCreate) (string, error) {
	args := strings.Split(m.Content[1:], " ")

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

// Returns a list of all stored custom commands
func ListCustomCommands(m *discordgo.MessageCreate) (*discordgo.MessageEmbed, error) {
	args := strings.Split(m.Content, " ")
	var pageNumber int

	// Verify that args are of correct format
	if len(args)-1 > 1 {
		return nil, errors.New("correct usage is !list <page number>")
	}

	// Get all commands from the database
	commands, err := database.GetAllCustomCommandNames()
	if err != nil {
		return nil, err
	}

	// Calculate how many pages there are
	pages := int(float64(len(commands) / 10))
	if len(commands)%10 != 0 {
		pages += 1 // We rounded down so need to add 1
	}

	// If no page argument is given, set to 1
	if len(args) == 1 {
		pageNumber = 1
	} else {
		pageNumber, err = strconv.Atoi(args[1])
		if err != nil {
			return nil, err
		}
	}

	// Avoid out of range errors
	if pageNumber > pages {
		return nil, fmt.Errorf("there are only %d pages of commands", pages)
	}

	// Populate the selected page
	if pageNumber < pages {
		commands = commands[(pageNumber-1)*10 : (pageNumber * 10)]
	} else {
		commands = commands[(pageNumber-1)*10:]
	}

	// Create the embed
	embed := discordgo.MessageEmbed{Title: "Custom commands", Description: fmt.Sprintf("Page %d of %d", pageNumber, pages)}

	// Add the fields to the embed
	for _, c := range commands {
		field := discordgo.MessageEmbedField{Name: "!" + c.Name, Value: c.Result, Inline: false}
		embed.Fields = append(embed.Fields, &field)
	}

	return &embed, nil
}
