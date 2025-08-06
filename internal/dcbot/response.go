package dcbot

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dark-person/remind-later-dc/internal/timeparse"
	"github.com/rs/zerolog/log"
)

// Delete given discord message with logging and error handle.
func deleteMessageWithLog(s *discordgo.Session, channelID string, messageID string) {
	err := s.ChannelMessageDelete(channelID, messageID)
	if err != nil {
		log.Error().Err(err).Msg("Error when delete message")
	}
	log.Debug().Str("messageID", messageID).Msg("Message deleted.")
}

// Send text message and added sent message to queue for cleanup purpose.
func (bm *BotManager) sendTextMessage(channelID string, text string) {
	msg, err := bm.session.ChannelMessageSend(channelID, text)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message.")
	}

	// Set message to queue
	bm.sentMsg = append(bm.sentMsg, msg)
}

// This function will be called (due to AddHandler above)
// every time a new message is created on any channel
// that the authenticated bot has access to.
func (bm *BotManager) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself, which is a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	log.Trace().Msgf("[%s] %s: %s", m.ChannelID, m.Author.Username, m.Content)

	// ----------------------------------------

	// Handle cleanup
	if m.Content == "!clean" || m.Content == "!clear" {
		// Delete user message
		deleteMessageWithLog(s, m.ChannelID, m.ID)

		// Clean outdated message
		bm.Cleanup()
	}

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

	log.Debug().Msgf("String detected: %s", message)

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

		// Send message
		bm.sendTextMessage(channelID, response)

		// Delete original message
		deleteMessageWithLog(s, m.ChannelID, m.ID)
	})
	log.Debug().Dur("duration", d).Msg("Delayed function set.")

	// Set emoji
	err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅") // 🔜
	if err != nil {
		log.Error().Err(err).Msg("Error when adding reaction")
	}
	log.Debug().Msg("Reaction added.")
}
