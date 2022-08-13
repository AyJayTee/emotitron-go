package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/AyJayTee/emotitron-go/components"
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
	commands = []string{"add", "remove"}

	// Add handlers
	s.AddHandler(messageCreate)

	// Store database connection
	components.StartDatabase()
	defer components.ShutdownDatabase()

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
			msg, err := invokeCommand(c, strings.Split(m.Content[len(c)+2:], " ")...)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}
			if msg != "" {
				s.ChannelMessageSend(m.ChannelID, msg)
			}
			return
		}
	}

	// Check for invoking custom command
	if strings.HasPrefix(m.Content, commandPrefix) {
		msg, err := components.GetCommand(strings.Split(m.Content[len(commandPrefix):], " "))
		if err != nil {
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
		components.PingDatabase()
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if m.Content == "createtable" {
		components.CreateTable()
	}
}

func invokeCommand(command string, args ...string) (string, error) {
	switch command {
	case "add":
		msg, err := components.AddCommand(args)
		if err != nil {
			return "", err
		}
		return msg, nil

	case "remove":
		msg, err := components.RemoveCommand(args)
		if err != nil {
			return "", err
		}
		return msg, nil
	}

	return "", nil
}
