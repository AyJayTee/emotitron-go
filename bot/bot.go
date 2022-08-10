package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Starts and returns a pointer to the bot session
func Start(token string) {
	s, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("Error creating session.", err)
		return
	}

	// Declare bot intents
	s.Identify.Intents = discordgo.IntentGuildMessages

	// Add handlers
	s.AddHandler(messageCreate)

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

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
