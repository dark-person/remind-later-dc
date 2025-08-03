package dcbot

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dark-person/remind-later-dc/internal/timeparse"
)

// Delete given discord message with logging and error handle.
func deleteMessageWithLog(s *discordgo.Session, channelID string, messageID string) {
	err := s.ChannelMessageDelete(channelID, messageID)
	if err != nil {
		fmt.Printf("[ERROR] Error delete message: %v", err)
	}
	fmt.Printf("[DEBUG] %s: Message %s deleted.\n", time.Now().Format("2006-01-02 15:04:05"), messageID)
}

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

	// Handle message with time string and mention
	if strings.HasPrefix(m.Content, s.State.User.Mention()) {
		bm.handleMention(s, m)
	}

	// ----------------------------------------
}

// Send message if message mention current bot as user.
func (bm *BotManager) handleMention(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Parse values from message object
	channelID := m.ChannelID
	author := m.Author

	// Remove mention from message
	message := strings.ReplaceAll(m.Content, s.State.User.Mention(), "")
	message = strings.TrimSpace(message)

	fmt.Printf("[DEBUG] %s : String detected: '%s'\n", time.Now().Format("2006-01-02 15:04:05"), message)

	// Split message into two part, time string and optional information
	splited := strings.SplitN(message, " ", 2)

	var timeStr, optMsg string
	if len(splited) == 2 {
		timeStr = splited[0]
		optMsg = splited[1]
	} else {
		timeStr = splited[0]
		optMsg = ""
	}

	// Extract time values from message
	hour, minutes, second, ok := timeparse.ExtractDuration(timeStr)

	// Ignore if message is not in correct format
	if !ok {
		return
	}

	// Get duration
	d := time.Hour*time.Duration(hour) + time.Minute*time.Duration(minutes) + time.Second*time.Duration(second)

	// Get time later
	now := time.Now()

	// Schedule job by time
	time.AfterFunc(d, func() {
		response := author.Mention() + ", reminder you that it is time for you to do some thing. (Requested at " + now.Format("2006-01-02 15:04:05") + ")"

		// Override response if optional message included
		if optMsg != "" {
			response = author.Mention() + ", reminder you that " + optMsg
		}

		msg, err := s.ChannelMessageSend(channelID, response)
		if err != nil {
			fmt.Printf("[ERROR] %s : %s %v\n", time.Now().Format("2006-01-02 15:04:05"), "Failed to send message.", err)
		}

		// Delete original message
		deleteMessageWithLog(s, m.ChannelID, m.ID)

		// Add handler for delete message after 24hrs
		time.AfterFunc(time.Duration(24)*time.Hour, func() {
			deleteMessageWithLog(s, msg.ChannelID, msg.ID)
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
