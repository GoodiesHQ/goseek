package main

import (
	"fmt"
	"os"
	"time"

	"github.com/goodieshq/goseek/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const PORT uint16 = 3000

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	authchecker := server.NewAuthCheck(
		map[string]string{
			"user": "passw0rd1",
		},
		[]string{
			"abc123",
		},
	)

	gs, err := server.NewGoSeek("./root", authchecker)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create a new GoSeek server instance")
	}

	log.Info().Str("addr", fmt.Sprintf(":%d", PORT)).Msg("listening")
	gs.Run(PORT)
}
