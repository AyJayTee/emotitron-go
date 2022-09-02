package components

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AyJayTee/emotitron-go/database"
	"github.com/bwmarrin/discordgo"
)

// Returns a map of commands provided by the custom commands component
func CustomCommands() map[string]func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return map[string]func(s *discordgo.Session, m *discordgo.MessageCreate){
		"add":    addCustomCommand,
		"remove": removeCustomCommand,
		"list":   listCustomCommands,
	}
}

// Adds a custom command to the database
func addCustomCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	hasAttachment := false
	args := strings.Split(m.Content, " ")

	// Check for attachments
	if len(m.Attachments) != 0 {
		hasAttachment = true
	}

	// Verify that the args are of correct format
	if hasAttachment {
		if len(args) != 2 {
			s.ChannelMessageSend(m.ChannelID, "Correct usage is !add [command name] [command value], or !add [command name] with an attachment")
			return
		}
	} else {
		if len(args) != 3 {
			s.ChannelMessageSend(m.ChannelID, "Correct usage is !add [command name] [command value], or !add [command name] with an attachment")
			return
		}
	}

	// Check that command does not already exist
	_, err := database.GetCustomCommandValue(args[1])
	if err == nil {
		s.ChannelMessageSend(m.ChannelID, "Cannot add a command that already exists")
		return
	}

	// Add the command to the database
	if hasAttachment {
		err := database.InsertCustomCommand(database.CustomCommand{Name: args[1], Result: m.Attachments[0].URL})
		if err != nil {
			log.Printf("Error %s when creating command with attachment", err)
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
	} else {
		err := database.InsertCustomCommand(database.CustomCommand{Name: args[1], Result: args[2]})
		if err != nil {
			log.Printf("Error %s when creating command", err)
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
	}

	s.ChannelMessageSend(m.ChannelID, "Command successfully added!")
}

// Removes a custom command from the database
func removeCustomCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")

	// Verify that the args are of correct format
	if len(args) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Correct usage is !remove [command name]")
		return
	}

	// Remove the command from the database
	err := database.RemoveCustomCommand(args[1])
	if err != nil {
		log.Printf("Error %s when removing command %s", err, args[1])
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Command successfully removed!")
}

// Returns a list of all stored custom commands
func listCustomCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")
	var pageNumber int

	// Verify that args are of correct format
	if len(args) > 2 {
		s.ChannelMessageSend(m.ChannelID, "Correct usage is !list <page number>")
		return
	}

	// Get all commands from the database
	commands, err := database.GetAllCustomCommandNames()
	if err != nil {
		log.Printf("Error %s when getting all commands", err)
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
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
			log.Printf("Error %s when converting string to int", err)
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
	}

	// Avoid out of range errors
	if pageNumber > pages {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("There are only %d pages of commands", pages))
		return
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

	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

// Fetches the value of a custom command from the database
func GetCustomCommand(m *discordgo.MessageCreate) (string, error) {
	args := strings.Split(m.Content, " ")

	// Verify that the args are of correct format
	if len(args) != 1 {
		return "", errors.New("wrong number of arguments")
	}

	// Fetch the command from the database
	result, err := database.GetCustomCommandValue(args[0][1:])
	if err != nil {
		return "", err
	}

	return result, nil
}
