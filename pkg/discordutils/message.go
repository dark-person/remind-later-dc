// Utils package for discord API.
package discordutils

import "github.com/bwmarrin/discordgo"

// Send message to specified channel.
func SendMsgToChannel(s *discordgo.Session, channelID string, msg string) error {
	// Send response message
	_, err := s.ChannelMessageSend(channelID, msg)
	if err != nil {
		return err // Could not send message.
	}

	return nil
}
