package database

import (
	"github.com/dgraph-io/badger"
	"encoding/json"

	"github.com/robbinhan/go-blockchain/types"
	"github.com/tendermint/tmlibs/log"
)

const (
	//LatestBlockKey
	LatestBlockKey = "LatestBlock"
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
	db     *badger.DB
	logger log.Logger
}

func NewState(db *badger.DB, logger log.Logger) *State {
	logger.With("module", "state")

	return &State{db: db, logger: logger}
}

func (state *State) Hash() {

}

/**
插入区块数据
 */
func (state *State) InsertBlock(block types.Block) error {
	blockBytes, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}
	err = state.db.Update(func(txn *badger.Txn) error {
		hash := block.Hash()
		state.logger.Debug("set block", "hash", hash)
		err := txn.Set(hash, blockBytes)
		err = txn.Set([]byte(LatestBlockKey), hash)
		return err
	})
	if err != nil {
		panic(err)
	}
	return nil
}

func (state *State) LatestBlock() (*types.Block, error) {
	var res []byte
	err := state.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(LatestBlockKey))
		if err == badger.ErrKeyNotFound {
			state.logger.Debug("not found val")
			return nil
		} else if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}

		state.logger.Debug("get LatestBlock hash", "hash", val)

		item, err = txn.Get(val)
		if err == badger.ErrKeyNotFound {
			state.logger.Info("not found val")
			return nil
		} else if err != nil {
			return err
		}
		val, err = item.Value()
		if err != nil {
			return err
		}

		state.logger.Debug("get LatestBlock", "block", val)

		res = val
		return nil
	})

	block := &types.Block{}
	if len(res) > 0 {
		err = json.Unmarshal(res, block)
		if err != nil {
			return block, err
		}
	}

	return block, err
}
