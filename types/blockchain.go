package types

import (
	"encoding/json"

	"github.com/tendermint/tmlibs/log"

	"github.com/dgraph-io/badger"
	"github.com/robbinhan/go-blockchain/database"
)

const (
	//LatestBlockKey
	LatestBlockKey = "LatestBlock"
)

type BlockChain struct {
	GenesisBlock Block
	ChainID      string
	state        *database.State
	logger       log.Logger
	CurrentBlock *Block
}

func NewBlockChain(state *database.State, logger log.Logger) *BlockChain {
	return &BlockChain{
		state:  state,
		logger: logger,
	}
}

// AddBlock 插入区块数据
func (blockchain *BlockChain) AddBlock(block Block) error {
	blockBytes, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}
	err = blockchain.state.DB.Update(func(txn *badger.Txn) error {
		hash := block.Hash()
		blockchain.logger.Debug("set block", "hash", hash)
		err := txn.Set(hash, blockBytes)
		err = txn.Set([]byte(LatestBlockKey), hash)
		return err
	})
	if err != nil {
		panic(err)
	}
	return nil
}

// LatestBlock 返回最新区块
func (blockchain *BlockChain) LatestBlock() (*Block, error) {
	var res []byte
	err := blockchain.state.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(LatestBlockKey))
		if err == badger.ErrKeyNotFound {
			blockchain.logger.Debug("not found val")
			return nil
		} else if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}

		blockchain.logger.Debug("get LatestBlock hash", "hash", val)

		item, err = txn.Get(val)
		if err == badger.ErrKeyNotFound {
			blockchain.logger.Info("not found val")
			return nil
		} else if err != nil {
			return err
		}
		val, err = item.Value()
		if err != nil {
			return err
		}

		blockchain.logger.Debug("get LatestBlock", "block", val)

		res = val
		return nil
	})

	block := &Block{}
	if len(res) > 0 {
		err = json.Unmarshal(res, block)
		if err != nil {
			return block, err
		}
	}

	blockchain.CurrentBlock = block

	return block, err
}
