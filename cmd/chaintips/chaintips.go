package main

import (
	"os"

	"github.com/nbitslabs/chaintips/api"
	"github.com/nbitslabs/chaintips/storage/sqlite"
	"github.com/nbitslabs/chaintips/tracker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: log.Output(zerolog.ConsoleWriter{Out: os.Stderr})})
}

func main() {
	db, err := sqlite.NewSqliteBackend()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database")
	}

	tracker := tracker.NewTracker(db)
	go tracker.Run()

	api := api.NewServer(db)
	api.Serve()
}
