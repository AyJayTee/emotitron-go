package components

import (
	"errors"
	"strings"

	"github.com/AyJayTee/emotitron-go/database"
	"github.com/bwmarrin/discordgo"
)

// Adds a new response to the database
func AddResponse(m *discordgo.MessageCreate) (string, error) {
	args := strings.Split(m.Content[13:], " ")
	// Check args are of correct format
	if len(args) < 2 {
		return "", errors.New("correct usage is !addresponse [trigger] [response]")
	}

	// Build the response
	response := database.Response{
		Trigger:  args[0],
		Response: strings.Join(args[1:], " "),
	}

	// Commit to the database
	err := database.InsertResponse(response)
	if err != nil {
		return "", err
	}

	return "Response successfully created!", nil
}

// Removes a response from the database
func RemoveResponse(m *discordgo.MessageCreate) (string, error) {
	args := strings.Split(m.Content[16:], " ")
	// Check args are of correct format
	if len(args) != 1 {
		return "", errors.New("correct usage is !removeresponse [trigger]")
	}

	// Remove from the database
	err := database.RemoveResponse(args[0])
	if err != nil {
		return "", err
	}

	return "Reponse successfully deleted!", nil
}

// Modifies the trigger of a response in the database
func ModifyTrigger(m *discordgo.MessageCreate) (string, error) {
	args := strings.Split(m.Content[15:], " ")
	// Check args are of correct format
	if len(args) != 2 {
		return "", errors.New("correct usage is !modifytrigger [trigger to modify] [new trigger]")
	}

	// Update the database
	err := database.UpdateResponseTrigger(args[0], args[1])
	if err != nil {
		return "", err
	}

	return "Response successfully updated!", nil
}

// Modifies the response of a response in the database
func MofifyResponse(m *discordgo.MessageCreate) (string, error) {
	args := strings.Split(m.Content[16:], " ")
	// Check args are of correct format
	if len(args) != 2 {
		return "", errors.New("correct usage is !modifyresponse [trigger] [new response]")
	}

	// Update the database
	err := database.UpdateResponseResponse(args[0], args[1])
	if err != nil {
		return "", err
	}

	return "Response successfully updated!", nil
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

// Lists all current repsonses
func ListResponses() (*discordgo.MessageEmbed, error) {
	// Get all responses
	responses, err := database.GetAllResponses()
	if err != nil {
		return nil, err
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

	return &embed, nil
}
