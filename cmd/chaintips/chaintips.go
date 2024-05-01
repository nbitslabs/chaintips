package main

import (
	"log"

	"github.com/nbitslabs/chaintips/storage/sqlite"
	"github.com/nbitslabs/chaintips/tracker"
)

func main() {
	db, err := sqlite.NewSqliteBackend()
	if err != nil {
		log.Fatal(err)
	}

	tracker := tracker.NewTracker(db)
	tracker.Run()
}
