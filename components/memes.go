package components

import (
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Randomly mixes up letters in the previous message and then posts the result
func ChrisTranslate(s *discordgo.Session, m *discordgo.MessageCreate) (string, error) {
	// Fetch the previous message
	history, err := s.ChannelMessages(m.ChannelID, 1, m.ID, "", "")
	if err != nil {
		return "", err
	}

	// Separate the message and randomly reorder letters
	modifier := 10 // 1 to 10 chance to mix up letters
	words := strings.Split(history[0].Content, " ")
	var newWords []string
	for _, w := range words {
		if rand.Intn(10) <= modifier {
			randomInt1 := rand.Intn(len(w))
			randomInt2 := rand.Intn(len(w))
			letter1 := w[randomInt1]
			letter2 := w[randomInt2]
			var newWord string
			for i, l := range w {
				if i == randomInt1 {
					newWord += string(letter2)
				} else if i == randomInt2 {
					newWord += string(letter1)
				} else {
					newWord += string(l)
				}
			}

			newWords = append(newWords, newWord)
		} else {
			newWords = append(newWords, w)
		}
	}
	msg := strings.Join(newWords, " ")

	return msg, nil
}
