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
	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Declare bot commands
	commands = []string{
		"help",                  // general.go
		"add", "remove", "list", // customcommands.go
		"christranslate",       // memes.go
		"remindme", "forgetme", // reminders.go
		"addresponse", "removeresponse", "modifytrigger", "modifyresponse", "listresponses", // responses.go
	}

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
		// Check for proper command
		for _, c := range commands {
			if words[0] == commandPrefix+c {
				err := invokeCommand(c, s, m)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, err.Error())
					return
				}
				return
			}
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

func invokeCommand(command string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	switch command {
	case "help":
		embed := components.Help()
		s.ChannelMessageSendEmbed(m.ChannelID, embed)

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

	case "christranslate":
		msg, err := components.ChrisTranslate(s, m)
		if err != nil {
			return err
		}
		s.ChannelMessageSend(m.ChannelID, msg)
		return nil

	case "remindme":
		msg, err := components.RemindMe(m)
		if err != nil {
			return err
		}
		s.ChannelMessageSend(m.ChannelID, msg)
		return nil

	case "forgetme":
		msg, err := components.ForgetMe(m)
		if err != nil {
			return err
		}
		s.ChannelMessageSend(m.ChannelID, msg)
		return nil

	case "addresponse":
		msg, err := components.AddResponse(m)
		if err != nil {
			return err
		}
		s.ChannelMessageSend(m.ChannelID, msg)
		return nil

	case "removeresponse":
		msg, err := components.RemoveResponse(m)
		if err != nil {
			return err
		}
		s.ChannelMessageSend(m.ChannelID, msg)
		return nil

	case "modifytrigger":
		msg, err := components.ModifyTrigger(m)
		if err != nil {
			return err
		}
		s.ChannelMessageSend(m.ChannelID, msg)
		return nil

	case "modifyresponse":
		msg, err := components.MofifyResponse(m)
		if err != nil {
			return err
		}
		s.ChannelMessageSend(m.ChannelID, msg)

	case "listresponses":
		msg, err := components.ListResponses()
		if err != nil {
			return err
		}
		s.ChannelMessageSendEmbed(m.ChannelID, msg)
	}

	return nil
}
