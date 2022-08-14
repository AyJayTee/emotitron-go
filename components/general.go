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

	// !remindme
	remindMeField := discordgo.MessageEmbedField{
		Name:   "!remindme [quantity] [minutes/hours/days/weeks/months] [text]",
		Value:  "PMs you the specified text after the specified time has elapsed",
		Inline: false,
	}
	embed.Fields = append(embed.Fields, &remindMeField)

	// !forgetme
	forgetMeField := discordgo.MessageEmbedField{
		Name:   "!forgetme",
		Value:  "Removes all of your current reminders",
		Inline: false,
	}
	embed.Fields = append(embed.Fields, &forgetMeField)

	// !addresponse
	addResponseField := discordgo.MessageEmbedField{
		Name:   "!addresponse [trigger] [response]",
		Value:  "Adds a response to the bot so that when a user writes sends the trigger, the bot responds with the response",
		Inline: false,
	}
	embed.Fields = append(embed.Fields, &addResponseField)

	// !removeresponse
	removeResponseField := discordgo.MessageEmbedField{
		Name:   "!removeresponse [trigger]",
		Value:  "Removes the specified response",
		Inline: false,
	}
	embed.Fields = append(embed.Fields, &removeResponseField)

	// !modifytrigger
	modifyTriggerField := discordgo.MessageEmbedField{
		Name:   "!modifytrigger [trigger to modify] [new trigger]",
		Value:  "Renames the trigger portion of a response",
		Inline: false,
	}
	embed.Fields = append(embed.Fields, &modifyTriggerField)

	// !modifyrepsonse
	modifyResponseField := discordgo.MessageEmbedField{
		Name:   "!modifyresponse [trigger] [new response]",
		Value:  "Modifies the response portion of a response",
		Inline: false,
	}
	embed.Fields = append(embed.Fields, &modifyResponseField)

	// !listresponses
	listResponsesField := discordgo.MessageEmbedField{
		Name:   "!listresponses",
		Value:  "Lists all currently stored responses",
		Inline: false,
	}
	embed.Fields = append(embed.Fields, &listResponsesField)

	return &embed
}
