package types

import (
	//tendertypes "github.com/tendermint/tendermint/types"
	"crypto/sha256"
	"github.com/robbinhan/go-blockchain/bazel-go-blockchain/external/go_sdk/src/encoding/json"
	"github.com/tendermint/go-amino"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/merkle"
	"golang.org/x/crypto/ripemd160"
)

type Block struct {
	Header
	Parent *Block
	Txs    []*Transaction
	Logs   []Log
}

type Header struct {
	ChainID  string `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Time     int64  `protobuf:"varint,3,opt,name=time,proto3" json:"time,omitempty"`
	NumTxs   int32  `protobuf:"varint,4,opt,name=num_txs,json=numTxs,proto3" json:"num_txs,omitempty"`
	Height   int64  `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	TotalTxs int64  `protobuf:"varint,5,opt,name=total_txs,json=totalTxs,proto3" json:"total_txs,omitempty"`
	// hashes
	LastBlockHash  []byte `protobuf:"bytes,6,opt,name=last_block_hash,json=lastBlockHash,proto3" json:"last_block_hash,omitempty"`
	ValidatorsHash []byte `protobuf:"bytes,7,opt,name=validators_hash,json=validatorsHash,proto3" json:"validators_hash,omitempty"`
	AppHash        []byte `protobuf:"bytes,8,opt,name=app_hash,json=appHash,proto3" json:"app_hash,omitempty"`
	// consensus
	//Proposer Validator `protobuf:"bytes,9,opt,name=proposer" json:"proposer"`
	MerkleRootHash []byte
}

// Hash returns the hash of the header.
// Returns nil if ValidatorHash is missing,
// since a Header is not valid unless there is
// a ValidaotrsHash (corresponding to the validator set).
func (h *Header) Hash() []byte {
	if h == nil {
		return nil
	}

	second := sha256.New()

	hm := map[string]interface{}{
		"ChainID":    h.ChainID,
		"Height":     h.Height,
		"Time":       h.Time,
		"NumTxs":     h.NumTxs,
		"TotalTxs":   h.TotalTxs,
		"Validators": h.ValidatorsHash,
		"App":        h.AppHash,
	}
	hmBytes, _ := json.Marshal(hm)

	second.Write(hmBytes)
	return second.Sum(nil)
	//return merkle.SimpleHashFromMap(map[string]merkle.Hasher{
	//	"ChainID":     aminoHasher(h.ChainID),
	//	"Height":      aminoHasher(h.Height),
	//	"Time":        aminoHasher(h.Time),
	//	"NumTxs":      aminoHasher(h.NumTxs),
	//	"TotalTxs":    aminoHasher(h.TotalTxs),
	//	//"LastBlockID": aminoHasher(h.LastBlockID),
	//	//"LastCommit":  aminoHasher(h.LastCommitHash),
	//	//"Data":        aminoHasher(h.DataHash),
	//	"Validators":  aminoHasher(h.ValidatorsHash),
	//	"App":         aminoHasher(h.AppHash),
	//	//"Consensus":   aminoHasher(h.ConsensusHash),
	//	//"Results":     aminoHasher(h.LastResultsHash),
	//	//"Evidence":    aminoHasher(h.EvidenceHash),
	//})
}

type hasher struct {
	item interface{}
}

var cdc = amino.NewCodec()

func (h hasher) Hash() []byte {
	hasher := ripemd160.New()
	if h.item != nil && !cmn.IsTypedNil(h.item) && !cmn.IsEmpty(h.item) {
		bz, err := cdc.MarshalBinaryBare(h.item)
		if err != nil {
			panic(err)
		}
		_, err = hasher.Write(bz)
		if err != nil {
			panic(err)
		}
	}
	return hasher.Sum(nil)

}

func aminoHash(item interface{}) []byte {
	h := hasher{item}
	return h.Hash()
}

func aminoHasher(item interface{}) merkle.Hasher {
	return hasher{item}
}

type Transaction struct {
	txID []byte
	data []byte
}

func NewTransaction(data []byte, block Block) *Transaction {
	tx := &Transaction{data: data}
	shaHasher := sha256.New()
	shaHasher.Write(data)
	shaHasher.Write(block.Hash())
	tx.txID = shaHasher.Sum(nil)
	return tx
}

type Log struct {
}
