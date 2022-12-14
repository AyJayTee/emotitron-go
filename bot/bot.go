package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	commands map[string]func(s *discordgo.Session, m *discordgo.MessageCreate)
	config   Config
)

type Config struct {
	Token  string `json:"bot_token"`
	Prefix string `json:"bot_prefix"`
}

// Starts and returns a pointer to the bot session
func Start() {
	// Open the config file
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Printf("Error %s opening config.json", err)
	}
	defer configFile.Close()

	// Read the config file
	byteValue, _ := ioutil.ReadAll(configFile)
	json.Unmarshal(byteValue, &config)

	// Extract the token and start the bot
	token := config.Token
	s, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("Error creating session.", err)
		return
	}

	// Declare bot intents
	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Initalize the commands map
	commands = make(map[string]func(s *discordgo.Session, m *discordgo.MessageCreate))

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
	if strings.HasPrefix(m.Content, config.Prefix) {
		checkForCommand(words, s, m)
	}

	// Check for response trigger
	checkForResponse(words, s, m)
}

func checkForCommand(words []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	// Guard clause for not beginning with prefix
	if !strings.HasPrefix(words[0], config.Prefix) {
		return
	}

	// Extract the command name
	command := words[0][1:]

	// Check for proper command
	if ok := invokeCommand(command, s, m); ok {
		return
	}

	// Ignore if more than one word is in the message
	if len(words) > 1 {
		return
	}

	// Check for custom command
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
}

func checkForResponse(words []string, s *discordgo.Session, m *discordgo.MessageCreate) {
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
