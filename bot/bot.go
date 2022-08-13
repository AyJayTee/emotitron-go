package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/AyJayTee/emotitron-go/components"
	"github.com/AyJayTee/emotitron-go/database"
	"github.com/bwmarrin/discordgo"
)

var (
	commands      []string
	commandPrefix string
)

// Starts and returns a pointer to the bot session
func Start() {
	token := os.Getenv("BOT_TOKEN")
	commandPrefix = os.Getenv("BOT_PREFIX")

	s, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("Error creating session.", err)
		return
	}

	// Declare bot intents
	s.Identify.Intents = discordgo.IntentGuildMessages

	// Declare bot commands
	commands = []string{"add", "remove", "list"}

	// Add handlers
	s.AddHandler(messageCreate)

	// Store database connection
	database.StartDatabase()
	defer database.ShutdownDatabase()

	// Open a connection
	err = s.Open()
	if err != nil {
		fmt.Println("Error opening connection.", err)
		return
	}
	defer s.Close()

	fmt.Println("Emotitron activated, press CTRL-C to stop.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check for invoking proper command
	for _, c := range commands {
		if strings.HasPrefix(m.Content, commandPrefix+c) {
			err := invokeCommand(c, s, m)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}
			return
		}
	}

	// Check for invoking custom command
	if strings.HasPrefix(m.Content, commandPrefix) {
		msg, err := components.GetCustomCommand(m)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				// Command does not exist so ignore this case
				return
			}
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		if msg != "" {
			s.ChannelMessageSend(m.ChannelID, msg)
		}
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
		database.PingDatabase()
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if m.Content == "createtable" {
		database.CreateTable()
	}
}

func invokeCommand(command string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	switch command {
	case "add":
		msg, err := components.AddCustomCommand(commands, m)
		if err != nil {
			return err
		}
		s.ChannelMessageSend(m.ChannelID, msg)
		return nil

	case "remove":
		msg, err := components.RemoveCustomCommand(m)
		if err != nil {
			return err
		}
		s.ChannelMessageSend(m.ChannelID, msg)
		return nil

	case "list":
		embed, err := components.ListCustomCommands()
		if err != nil {
			return err
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return nil
	}

	return nil
}
