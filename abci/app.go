package abci

import (
	"os"

	"github.com/dgraph-io/badger"
	"github.com/robbinhan/go-blockchain/database"
	"github.com/robbinhan/go-blockchain/types"
	abcitypes "github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/log"
)

type Application struct {
	abcitypes.BaseApplication
	logger     log.Logger
	block      types.Block
	validators []abcitypes.Validator
	blockchain *types.BlockChain
}

func NewApplication(db *badger.DB) *Application {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	state := database.NewState(db, logger)
	blockchain := types.NewBlockChain(state, logger)

	return &Application{logger: logger, blockchain: blockchain}
}

/**
Info Connection  start
*/
func (app *Application) Info(req abcitypes.RequestInfo) (resInfo abcitypes.ResponseInfo) {
	app.logger.Debug("Info", "version", req.GetVersion(), "string", req.String())
	// 检查db，查询当前height，返回不同的数据
	currentBlock, err := app.blockchain.LatestBlock()
	if err != nil {
		panic(err)
	}
	height := currentBlock.Height
	hash := currentBlock.AppHash

	app.logger.Debug("Info", "current_height", height) // nolint: errcheck

	if height == 0 {
		return abcitypes.ResponseInfo{
			Data:             "ABCI",
			LastBlockHeight:  height,
			LastBlockAppHash: []byte{},
		}
	}

	return abcitypes.ResponseInfo{
		Data:             "ABCIEthereum",
		LastBlockHeight:  height,
		LastBlockAppHash: hash[:],
	}
}

func (app *Application) SetOptionn(req abcitypes.RequestSetOption) (resOption abcitypes.Response_SetOption) {
	app.logger.Debug("SetOption", "key", req.Key, "value", req.Value)
	return abcitypes.Response_SetOption{}
}

func (app *Application) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	app.logger.Debug("Query")
	return abcitypes.ResponseQuery{Log: reqQuery.Path}
}

/**
Info Connection  end
*/

/**
Mempool Connection
*/
func (app *Application) CheckTx(tx []byte) abcitypes.ResponseCheckTx {
	app.logger.Debug("CheckTx")
	return abcitypes.ResponseCheckTx{Code: 0, Data: tx, Fee: cmn.KI64Pair{Key: []byte("key"), Value: 1}}
}

/**
Consensus Connection start
*/
func (app *Application) InitChain(req abcitypes.RequestInitChain) (resInit abcitypes.ResponseInitChain) {
	app.logger.Debug("InitChain", "validator", req.GetValidators(), "state", req.GetAppStateBytes())

	app.validators = req.GetValidators()
	resp := abcitypes.ResponseInitChain{
		Validators: req.Validators, // 不改变validator
	}
	return resp
}

// Track the block hash and header information
func (app *Application) BeginBlock(params abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.logger.Debug("BeginBlock", "hash", params.GetHash(), "header", params.Header, "Validator", params.GetValidators())
	app.block.Header.ChainID = params.GetHeader().ChainID
	app.block.Header.Time = params.GetHeader().Time
	app.block.Header.NumTxs = params.GetHeader().NumTxs
	app.block.Header.TotalTxs = params.GetHeader().TotalTxs
	app.block.Header.LastBlockHash = params.GetHeader().LastBlockHash
	app.block.Header.ValidatorsHash = params.GetHeader().ValidatorsHash
	app.block.Header.AppHash = params.GetHeader().AppHash
	return abcitypes.ResponseBeginBlock{}
}

// tx is either "key=value" or just arbitrary bytes
func (app *Application) DeliverTx(tx []byte) abcitypes.ResponseDeliverTx {
	app.logger.Debug("DeliverTx")
	app.logger.Debug("print tx", "tx", tx)

	app.block.Txs = append(app.block.Txs, types.NewTransaction(tx, app.block))
	return abcitypes.ResponseDeliverTx{Code: 0}
}

func (app *Application) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	app.logger.Debug("EndBlock", "height", req.Height)
	app.block.Height = req.Height

	resp := abcitypes.ResponseEndBlock{}
	return resp
}

func (app *Application) Commit() abcitypes.ResponseCommit {
	app.logger.Debug("Commit")
	appHash := make([]byte, 8)
	rootHash := types.RootHash(app.block)
	app.block.Header.MerkleRootHash = rootHash
	app.block.Header.LastBlockHash = app.blockchain.CurrentBlock.Hash()
	app.block.Parent = app.blockchain.CurrentBlock
	app.blockchain.AddBlock(app.block)
	return abcitypes.ResponseCommit{Data: appHash}
}

/**
Consensus Connection end
*/
