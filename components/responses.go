package components

import (
	"log"
	"strings"

	"github.com/AyJayTee/emotitron-go/database"
	"github.com/bwmarrin/discordgo"
)

// Returns a map of commands provided by the responses component
func Responses() map[string]func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return map[string]func(s *discordgo.Session, m *discordgo.MessageCreate){
		"addresponse":    addResponse,
		"removeresponse": removeResponse,
		"modifytrigger":  modifyTrigger,
		"modifyresponse": modifyResponse,
		"listresponses":  listResponses,
	}
}

// Adds a new response to the database
func addResponse(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")
	// Check args are of correct format
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Correct usage is !addresponse [trigger] [response]")
		return
	}

	// Build the response
	response := database.Response{
		Trigger:  args[1],
		Response: strings.Join(args[2:], " "),
	}

	// Commit to the database
	err := database.InsertResponse(response)
	if err != nil {
		log.Printf("Error %s when adding response to database", err)
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Response successfully created!")
}

// Removes a response from the database
func removeResponse(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")
	// Check args are of correct format
	if len(args) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Correct usage is !removeresponse [trigger]")
		return
	}

	// Remove from the database
	err := database.RemoveResponse(args[1])
	if err != nil {
		log.Printf("Error %s when removing response from the database", err)
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Response successfully deleted!")
}

// Modifies the trigger of a response in the database
func modifyTrigger(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")
	// Check args are of correct format
	if len(args) != 3 {
		s.ChannelMessageSend(m.ChannelID, "Correct usage is !modifytrigger [trigger to modify] [new trigger]")
		return
	}

	// Update the database
	err := database.UpdateResponseTrigger(args[1], args[2])
	if err != nil {
		log.Printf("Error %s when updating the database", err)
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Trigger successfully updated!")
}

// Modifies the response of a response in the database
func modifyResponse(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")
	// Check args are of correct format
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Correct usage is !modifyresponse [trigger] [new response]")
		return
	}

	// Update the database
	err := database.UpdateResponseResponse(args[1], strings.Join(args[2:], " "))
	if err != nil {
		log.Printf("Error %s when updating the database", err)
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Response successfully updated!")
}

// Lists all current repsonses
func listResponses(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Get all responses
	responses, err := database.GetAllResponses()
	if err != nil {
		log.Printf("Error %s when fetching responses", err)
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	// Build the embed
	embed := discordgo.MessageEmbed{Title: "All current responses", Description: ""}
	for _, r := range responses {
		embedField := discordgo.MessageEmbedField{
			Name:   r.Trigger,
			Value:  r.Response,
			Inline: false,
		}
		embed.Fields = append(embed.Fields, &embedField)
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

// Checks if a trigger exists
func CheckForResponseTrigger(trigger string) bool {
	_, err := database.GetResponse(trigger)
	return err == nil
}

// Gets the value of a response
func GetResponseValue(trigger string) (database.Response, error) {
	response, err := database.GetResponse(trigger)
	if err != nil {
		return database.Response{}, err
	}

	return response, nil
}
