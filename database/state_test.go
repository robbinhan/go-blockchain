package database

import (
	"github.com/dgraph-io/badger"
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	db, _ := NewDB("testdata")

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("test"), []byte("hello"))
		return err
	})
	if err != nil {
		panic(err)
	}
}

func TestLatestBlock(t *testing.T) {
	var res []byte
	db, _ := NewDB("testdata")
	db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("test"))
		if err == badger.ErrKeyNotFound {
			fmt.Println("not found val")
			return nil
		} else if err != nil{
			fmt.Println("err", err)
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}
		fmt.Println("val", val)
		res = val
		return nil
	})
}
