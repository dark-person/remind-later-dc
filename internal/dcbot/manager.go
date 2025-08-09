// Package dcbot is responsible to manage discord bot behavior.
package dcbot

import (
	"fmt"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/dark-person/remind-later-dc/internal/config"
	"github.com/rs/zerolog/log"
)

// Manager for control static functions reference of this discord package.
//
// Before using this manager, you should call Init function like this:
//
//	m := dcbot.NewManager()
//	m.Init(someConfig)
//
// Otherwise, this bot manager will never work properly.
type BotManager struct {
	initialized bool               // Only true when this manager is initialized
	session     *discordgo.Session // Discord session that designed for notification

	Channel string // Channel ID to listen

	sentMsg []*discordgo.Message // Sent discord message, for cleanup purpose
}

// Create a new empty discord bot manager.
func NewManager() *BotManager {
	return &BotManager{
		initialized: false,
		session:     nil,
		Channel:     "",
		sentMsg:     make([]*discordgo.Message, 0),
	}
}

// Init this bot manager with given configuration,
// which also validate the configuration is able to run or not.
func (bm *BotManager) Init(cfg *config.DiscordConfig) error {
	// Perform validation of the configuration
	if cfg.Token == "" || cfg.ListenedChannel == "" {
		return fmt.Errorf("discord token or channel ID not set")
	}

	// Set channel from config
	bm.Channel = cfg.ListenedChannel

	var err error

	// Create a discord connection session
	bm.session, err = discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return fmt.Errorf("failed to create discord bot: %v", err)
	}

	// Add interaction listener for message
	bm.session.Identify.Intents |= discordgo.IntentMessageContent
	bm.session.AddHandler(bm.messageCreate)

	err = bm.session.Open()
	if err != nil {
		return fmt.Errorf("failed to open discord connection: %v", err)
	}

	bm.initialized = true
	return nil
}

// Clean all message that bot sent.
func (bm *BotManager) Cleanup() error {
	// Create map due to bulk delete message API require channel id
	m := make(map[string][]string, 0)

	// Loop send message
	for _, msg := range bm.sentMsg {
		m[msg.ChannelID] = append(m[msg.ChannelID], msg.ID)
	}

	// Loop map
	for k, v := range m {
		// Prepare iterator due to discord API bulk delete limit=100
		iter := slices.Chunk(v, 100)

		for batch := range iter {
			err := bm.session.ChannelMessagesBulkDelete(k, batch)
			if err != nil {
				return err
			}
		}
	}

	log.Debug().Msg("All sent message cleanup.")
	return nil
}

// Close session of the bot, with all message cleanup
func (bm *BotManager) CloseWithCleanup() error {
	err := bm.Cleanup()
	if err != nil {
		return err
	}

	return bm.session.Close()
}
