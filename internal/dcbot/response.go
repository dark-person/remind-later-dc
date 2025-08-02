package dcbot

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dark-person/remind-later-dc/internal/timeparse"
)

// This function will be called (due to AddHandler above)
// every time a new message is created on any channel
// that the authenticated bot has access to.
func (bm *BotManager) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself, which is a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Printf("%s [%s] %s: %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		m.ChannelID, m.Author.Username, m.Content)

	// ----------------------------------------

	// Check if bot is mentioned by '@' method
	bm.replyIfMentionBot(s, m)

	// ----------------------------------------
}

// Send message if message mention current bot as user.
func (bm *BotManager) replyIfMentionBot(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Parse values from message object
	channelID := m.ChannelID
	author := m.Author
	message := m.Content

	// Ignore if message not contain bot mention
	if !strings.HasPrefix(message, s.State.User.Mention()) {
		return
	}

	// Remove mention from message
	str := strings.ReplaceAll(message, s.State.User.Mention(), "")
	str = strings.TrimSpace(str)

	fmt.Printf("[DEBUG] %s : String detected: '%s'\n", time.Now().Format("2006-01-02 15:04:05"), str)

	// Split message into two part, time string and optional information
	splited := strings.SplitN(str, " ", 2)

	timeStr := splited[0]
	optMsg := splited[1]

	// Extract time values from message
	hour, minutes, second, ok := timeparse.ExtractDuration(timeStr)

	// Ignore if message is not in correct format
	if !ok {
		return
	}

	// Get duration
	d := time.Hour*time.Duration(hour) + time.Minute*time.Duration(minutes) + time.Second*time.Duration(second)

	// Get time later
	t := time.Now().Add(d)

	// Schedule job by time
	time.AfterFunc(d, func() {
		response := author.Mention() + ", reminder send at " + t.Format("2006-01-02 15:04:05")

		// Override response if optional message included
		if optMsg != "" {
			response = author.Mention() + ", reminder you that " + optMsg
		}

		msg, err := s.ChannelMessageSend(channelID, response)
		if err != nil {
			fmt.Printf("[ERROR] %s : %s %v\n", time.Now().Format("2006-01-02 15:04:05"), "Failed to send message.", err)
		}

		// Delete original message
		err = s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			fmt.Printf("[ERROR] Error delete message: %v", err)
		}
		fmt.Printf("[DEBUG] %s: Original message deleted.\n", time.Now().Format("2006-01-02 15:04:05"))

		// Add handler for delete message after 24hrs
		time.AfterFunc(time.Duration(24)*time.Hour, func() {
			err = s.ChannelMessageDelete(msg.ChannelID, msg.ID)
			if err != nil {
				fmt.Printf("[ERROR] Error delete message: %v", err)
			}
			fmt.Printf("[DEBUG] %s: Message clearup completed.\n", time.Now().Format("2006-01-02 15:04:05"))
		})
	})
	fmt.Printf("[DEBUG] %s : Delayed function set: %v\n", time.Now().Format("2006-01-02 15:04:05"), d)

	// Set emoji
	err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅") // 🔜
	if err != nil {
		fmt.Printf("[ERROR] Error adding reaction: %v", err)
	}
	fmt.Printf("[DEBUG] %s :Reaction added.\n", time.Now().Format("2006-01-02 15:04:05"))
}
