package abci

import (
	abcitypes "github.com/tendermint/abci/types"
	"os"
	"github.com/tendermint/tmlibs/log"
	"github.com/dgraph-io/badger"
	"github.com/robbinhan/go-blockchain/database"
	"github.com/robbinhan/go-blockchain/types"
)

type Application struct {
	abcitypes.BaseApplication
	blockHeader  abcitypes.Header
	state        *database.State
	logger       log.Logger
	block        types.Block
	validators   []abcitypes.Validator
	currentBlock *types.Block
}

func NewApplication(db *badger.DB) *Application {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	state := database.NewState(db, logger)

	block, err := state.LatestBlock()
	if err != nil {
		panic(err)
	}

	logger.Info("current block", "height", block.Height, "block", block)

	return &Application{state: state, logger: logger, currentBlock: block}
}

/**
Info Connection  start
 */
func (app *Application) Info(req abcitypes.RequestInfo) (resInfo abcitypes.ResponseInfo) {
	app.logger.Info("Info", "version", req.GetVersion(), "string", req.String())
	// 检查db，查询当前height，返回不同的数据
	//blockchain := app.backend.Ethereum().BlockChain()
	//currentBlock := blockchain.CurrentBlock()
	height := app.currentBlock.Height
	hash := app.currentBlock.Hash()

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
	app.logger.Info("SetOption", "key", req.Key, "value", req.Value)
	return abcitypes.Response_SetOption{}
}

func (app *Application) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	app.logger.Info("Query")
	return abcitypes.ResponseQuery{Log: reqQuery.Path}
}

/**
Info Connection  end
 */

/**
Mempool Connection
 */
func (app *Application) CheckTx(tx []byte) abcitypes.ResponseCheckTx {
	app.logger.Info("CheckTx")
	return abcitypes.ResponseCheckTx{Code: 1}
}

/**
Consensus Connection start
 */
func (app *Application) InitChain(req abcitypes.RequestInitChain) (resInit abcitypes.ResponseInitChain) {
	app.logger.Info("InitChain", "validator", req.GetValidators(), "state", req.GetAppStateBytes())

	app.validators = req.GetValidators()
	resp := abcitypes.ResponseInitChain{
		Validators: req.Validators, // 不改变validator
	}
	return resp
}

// Track the block hash and header information
func (app *Application) BeginBlock(params abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.logger.Info("BeginBlock", "hash", params.GetHash(), "header", params.Header, "Validator", params.GetValidators())
	app.block.Header.ChainID = params.GetHeader().ChainID
	app.block.Header.Time = params.GetHeader().Time
	app.block.Header.NumTxs = params.GetHeader().NumTxs
	app.block.Header.TotalTxs = params.GetHeader().TotalTxs
	app.block.Header.LastBlockHash = params.GetHeader().LastBlockHash
	app.block.Header.ValidatorsHash = params.GetHeader().ValidatorsHash
	app.block.Header.AppHash = params.GetHeader().AppHash

	app.blockHeader = params.GetHeader()

	return abcitypes.ResponseBeginBlock{}
}

// tx is either "key=value" or just arbitrary bytes
func (app *Application) DeliverTx(tx []byte) abcitypes.ResponseDeliverTx {
	app.logger.Info("DeliverTx")

	//parts := strings.Split(string(tx), "=")
	//if len(parts) == 2 {
	//	app.state.Set([]byte(parts[0]), []byte(parts[1]))
	//} else {
	//	app.state.Set(tx, tx)
	//}
	//return abcitypes.OK

	return abcitypes.ResponseDeliverTx{Code: 1}
}

func (app *Application) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	app.logger.Info("EndBlock", "height", req.Height)
	app.block.Height = req.Height

	resp := abcitypes.ResponseEndBlock{}
	return resp
}

func (app *Application) Commit() abcitypes.ResponseCommit {
	app.logger.Info("Commit")
	//hash := app.state.Hash()
	//return abcitypes.NewResultOK(hash, "")
	appHash := make([]byte, 8)
	//binary.PutVarint(appHash, app.state.Size)
	//app.state.AppHash = appHash
	//app.state.Height += 1
	rootHash := types.RootHash(app.block)
	app.block.Header.MerkleRootHash = rootHash
	app.state.InsertBlock(app.block)
	return abcitypes.ResponseCommit{Data: appHash}
}

/**
Consensus Connection end
 */
