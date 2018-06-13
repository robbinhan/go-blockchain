package main

import (
	"github.com/tendermint/abci/server"
	"os"
	"github.com/tendermint/tmlibs/log"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/abci/types"
	"github.com/robbinhan/go-blockchain/abci"
	"github.com/robbinhan/go-blockchain/database"
	"github.com/dgraph-io/badger"
)

func main() {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))

	// 数据库
	dataDir := "./data"
	db, err := database.NewDB(dataDir)
	defer db.Close()

	// 启动abci
	err = StartABCIServer("tcp://0.0.0.0:46658", "socket", db)
	if err != nil {
		logger.Error("Start abci server fail", err)
	}
}

func StartABCIServer(protoAddr, transport string, db *badger.DB) error {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))

	var app types.Application
	app = abci.NewApplication(db)

	//app = kvstore.NewPersistentKVStoreApplication("./")
	//app.(*kvstore.PersistentKVStoreApplication).SetLogger(logger.With("module", "kvstore"))

	srv, err := server.NewServer(protoAddr, transport, app)
	if err != nil {
		return err
	}
	srv.SetLogger(logger.With("module", "abci-server"))
	if err := srv.Start(); err != nil {
		return err
	}

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})
	return nil
}
