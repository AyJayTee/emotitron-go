package components

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/AyJayTee/emotitron-go/database"
	"github.com/bwmarrin/discordgo"
)

var (
	units = map[string]int{
		"minute": 60,
		"hour":   3600,
		"day":    86400,
		"week":   604800,
		"month":  2592000,
	}
)

// Returns a map of commands provided by the reminders component
func Reminders() map[string]func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return map[string]func(s *discordgo.Session, m *discordgo.MessageCreate){
		"remindme": remindMe,
		"forgetme": forgetMe,
	}
}

// Creates a reminder
func remindMe(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")

	// Check args are of the correct format
	if len(args) < 4 {
		s.ChannelMessageSend(m.ChannelID, "Correct usage is !remindme [quantity] [minutes/hours/days/weeks/months] [text]")
		return
	}

	// Compute the future time
	var future int
	quantity, err := strconv.ParseInt(args[1], 0, 32)
	if err != nil {
		log.Printf("Error %s when computing future time", err)
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	unit := strings.TrimSuffix(args[2], "s") // Trim the s from the unit definition
	if val, ok := units[unit]; ok {
		future = int(time.Now().Unix()) + (int(quantity) * val)
	} else {
		log.Println("Error when computing future time")
		s.ChannelMessageSend(m.ChannelID, "Unit not recognised, use one of minutes/hours/days/weeks/months")
		return
	}

	// Build the message string from remaining args
	text := strings.Join(args[3:], " ")

	// Build the reminder
	reminder := database.Reminder{
		UserID:    m.Author.ID,
		Future:    future,
		Text:      text,
		Completed: false,
	}

	// Commit reminder to the database
	err = database.InsertReminder(reminder)
	if err != nil {
		log.Printf("Error %s when committing reminder to database", err)
		s.ChannelMessageSend(m.ChannelID, "Error adding reminder to database")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Reminder successfully created!")
}

// Deletes all reminders
func forgetMe(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := database.RemoveReminders(m.Author.ID)
	if err != nil {
		log.Printf("Error %s when removing reminder", err)
		s.ChannelMessageSend(m.ChannelID, "Error when removing reminder")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "All of your reminders have been deleted!")
}

// Worker to check for due reminders
func CheckReminders(r chan<- database.Reminder) {
	for {
		// Find all expired reminders
		expiredReminders, err := database.GetExpiredReminders(time.Now().Unix())
		if err != nil {
			log.Printf("Error %s when checking reminders", err)
		}

		for _, reminder := range expiredReminders {
			// Send the reminder to the sender
			r <- reminder

			// Delete the expired reminder from the database
			err = database.RemoveRemdinder(reminder.Id)
			if err != nil {
				log.Printf("Error %s when clearing reminder", err)
			}
		}

		// Want to run every 5 seconds
		time.Sleep(5 * time.Second)
	}
}

// Worker to send out reminders recieved from the checker
func SendReminders(s *discordgo.Session, r <-chan database.Reminder) {
	for {
		reminder := <-r
		dmChannel, err := s.UserChannelCreate(reminder.UserID)
		if err != nil {
			log.Println("Error creating private channel with:", reminder.UserID)
		}
		// Build the message
		msg := "You asked me to remind you this: \n" + reminder.Text
		s.ChannelMessageSend(dmChannel.ID, msg)
		log.Println("Private messgae sent to", reminder.UserID)
	}
}
