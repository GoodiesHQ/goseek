package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goodieshq/goseek/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const PORT uint16 = 3000
const CFG_DEFAULT = "/app/config.yml"

func run() {
	// get the config path from the environment variable
	cfgPath := os.Getenv("GOSEEK_CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = CFG_DEFAULT
	}

	// create the server instance
	gs, err := server.NewGoSeek(cfgPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create a new GoSeek server instance")
	}

	// listen for SIGUP to reload the API keys
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)

	go func() {
		for {
			sig := <-sigs
			if sig == syscall.SIGHUP {
				log.Warn().Msg("reloading API keys")
				if err := gs.ReloadApiKeys(); err != nil {
					log.Error().Err(err).Msg("failed to reload API keys")
				}
			}
		}
	}()

	// run goseek
	log.Info().Int("pid", os.Getpid()).Str("addr", fmt.Sprintf(":%d", PORT)).Msg("goseek is running")
	gs.Run()
}

func main() {
	// initialize the logger
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	})
	run()
}
