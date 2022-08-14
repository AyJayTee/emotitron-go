package components

import "github.com/bwmarrin/discordgo"

func Help() *discordgo.MessageEmbed {
	// Create the embed
	embed := discordgo.MessageEmbed{Title: "Commands available", Description: ""}

	// !add
	addField := discordgo.MessageEmbedField{
		Name:   "!add [command name] [command result] or !add [command name] with an attachment",
		Value:  "Adds a new custom command and saves it to the database for future use.",
		Inline: false,
	}
	embed.Fields = append(embed.Fields, &addField)

	// !remove
	removeField := discordgo.MessageEmbedField{
		Name:   "!remove [command name]",
		Value:  "Removes a custom command from the database.",
		Inline: false,
	}
	embed.Fields = append(embed.Fields, &removeField)

	// !list
	listField := discordgo.MessageEmbedField{
		Name:   "!list",
		Value:  "Lists all custom commands currently stored in the database.",
		Inline: false,
	}
	embed.Fields = append(embed.Fields, &listField)

	// !christranslate
	christTranslateField := discordgo.MessageEmbedField{
		Name:   "!christranslate",
		Value:  "Utilises the dyslexia engine to translate the previous message into 'Chris speak'",
		Inline: false,
	}
	embed.Fields = append(embed.Fields, &christTranslateField)

	return &embed
}
