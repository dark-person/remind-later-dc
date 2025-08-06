package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dark-person/remind-later-dc/internal/config"
	"github.com/dark-person/remind-later-dc/internal/dcbot"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cfg *config.DiscordConfig

func main() {
	var err error
	// Load config
	cfg, err = config.LoadYaml("config.yaml")
	if err != nil {
		panic(err) // Program will never run properly when config not loaded
	}

	// Setup logger
	setupLogger(zerolog.TraceLevel)
	log.Info().Str("token", cfg.Token).Str("channels", cfg.ListenedChannel).Msg("Config loaded.")

	// Setup bot
	bot := dcbot.NewManager()
	err = bot.Init(cfg)
	if err != nil {
		panic(err)
	}

	// Wait for a termination signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Clean up
	err = bot.CloseWithCleanup()
	if err != nil {
		panic(err)
	}
}
