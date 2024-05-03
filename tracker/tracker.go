package tracker

import (
	"sync"

	"github.com/nbitslabs/chaintips/storage"
	"github.com/rs/zerolog/log"
)

var logger = log.With().Str("module", "tracker").Logger()

var syncWg sync.WaitGroup

type Tracker struct {
	db storage.Storage
}

func NewTracker(db storage.Storage) *Tracker {
	return &Tracker{db: db}
}

func (t *Tracker) Run() {
	chains, err := t.db.GetChains()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get chains")
		return
	}

	for _, chain := range chains {
		go t.indexBlocks(chain)
		go t.trackTips(chain)
		syncWg.Add(2)
	}

	syncWg.Wait()
}
