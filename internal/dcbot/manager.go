package dcbot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dark-person/remind-later-dc/internal/config"
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
func (bm *BotManager) Cleanup() {
	// Loop all message ID for deletion
	for _, msg := range bm.sentMsg {
		deleteMessageWithLog(bm.session, msg.ChannelID, msg.ID)
	}

	fmt.Printf("[DEBUG] %s: All sent message cleanup.\n", time.Now().Format("2006-01-02 15:04:05"))
}

// Close session of the bot, with all message cleanup
func (bm *BotManager) CloseWithCleanup() error {
	bm.Cleanup()
	return bm.session.Close()
}
