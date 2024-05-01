package tracker

import "github.com/nbitslabs/chaintips/storage"

type Tracker struct {
	db storage.Storage
}

func NewTracker(db storage.Storage) *Tracker {
	return &Tracker{db: db}
}

func (t *Tracker) Run() {
}
