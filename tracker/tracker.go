package tracker

import (
	"sync"

	"github.com/nbitslabs/chaintips/storage"
	"github.com/rs/zerolog/log"
)

var logger = log.With().Str("module", "tracker").Logger()

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

	syncWg := sync.WaitGroup{}

	for _, chain := range chains {
		syncWg.Add(1)
		go func() {
			defer syncWg.Done()
			t.indexBlocks(chain)
		}()

		syncWg.Add(1)
		go func() {
			defer syncWg.Done()
			t.trackTips(chain)
		}()

		syncWg.Add(1)
		go func() {
			defer syncWg.Done()
			t.linkChainTips(chain)
		}()
	}

	syncWg.Wait()
}
