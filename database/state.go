package database

import (
	"github.com/dgraph-io/badger"

	"github.com/tendermint/tmlibs/log"
)

func NewDB(dataDir string) (db *badger.DB, err error) {
	opts := badger.DefaultOptions
	opts.Dir = dataDir
	opts.ValueDir = dataDir
	db, err = badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type State struct {
	DB     *badger.DB
	logger log.Logger
}

func NewState(db *badger.DB, logger log.Logger) *State {
	logger.With("module", "state")
	return &State{DB: db, logger: logger}
}

func (state *State) Hash() {

}
