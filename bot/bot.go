package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/AyJayTee/emotitron-go/components"
	"github.com/AyJayTee/emotitron-go/database"
	"github.com/bwmarrin/discordgo"
)

var (
	commands      map[string]func(s *discordgo.Session, m *discordgo.MessageCreate)
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
	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Declare bot commands
	addCommands(components.CustomCommands()) // customcommands.go
	addCommands(components.General())        // general.go
	addCommands(components.Memes())          // memes.go
	addCommands(components.Reminders())      // reminders.go
	addCommands(components.Responses())      // responses.go

	// Add handlers
	s.AddHandler(messageCreate)

	// Open a connection
	err = s.Open()
	if err != nil {
		fmt.Println("Error opening connection.", err)
		return
	}
	defer s.Close()

	// Start database connection
	database.StartDatabase()
	defer database.ShutdownDatabase()

	// Start the reminder workers
	reminders := make(chan database.Reminder, 1)
	go components.CheckReminders(reminders)
	go components.SendReminders(s, reminders)

	fmt.Println("Emotitron activated, press CTRL-C to stop.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	words := strings.Split(m.Content, " ")

	// Ignore messages from the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check for invoking command
	if strings.HasPrefix(m.Content, commandPrefix) {
		// Extract the command name
		command := strings.Split(m.Content[1:], " ")[0]

		// Check for proper command
		if ok := invokeCommand(command, s, m); ok {
			return
		}

		// Check for custom command
		if len(words) == 1 {
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
	}

	// Check for response trigger
	for _, w := range words {
		if components.CheckForResponseTrigger(w) {
			response, err := components.GetResponseValue(w)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				log.Println(err.Error())
				return
			}
			s.ChannelMessageSend(m.ChannelID, response.Response)
			return
		}
	}
}

func invokeCommand(command string, s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if command, ok := commands[command]; ok {
		command(s, m)
		return true
	}
	return false
}

func addCommands(commandsToAdd map[string]func(s *discordgo.Session, m *discordgo.MessageCreate)) {
	for k, v := range commandsToAdd {
		commands[k] = v
	}
}
