package types

import (
	"fmt"
	"strconv"
	"testing"
)

func TestSplitNodeGroup(t *testing.T) {
	n := 5
	txs := mockTxs(n)

	var leafs []*Node

	for _, tx := range txs {
		l := NewNode(tx.txID)
		leafs = append(leafs, l)
	}

	groups := SplitNodeGroup(leafs)

	fmt.Println("groups", groups)

}

func TestRootHash(t *testing.T) {

	n := 4
	txs := mockTxs(n)

	block := Block{}
	block.Txs = txs

	rootHash := RootHash(block)

	fmt.Println("rootHash", rootHash)

	//getTracPath(txs, searchNode)

	//spv client
	//txHash -> req full node get merkle path

	// full node
	// get block by txHash
	// get merkle root hash by block header
	// get txs data by block body
	// get merkle path by cal all node hash

}
func TestProofPath(t *testing.T) {

	n := 4
	txs := mockTxs(n)

	block := Block{} // 需要证明的交易所在的block
	block.Txs = txs

	txIndex := 1 // 需要证明的交易索引

	var leafs []*Node

	for _, tx := range block.Txs {
		l := NewNode(tx.txID)
		leafs = append(leafs, l)
	}

	rootNode, allNodes := buildTree(leafs)
	allNodes = append(allNodes, leafs...)
	fmt.Println("rootNode", rootNode.Hash)
	for _, node := range allNodes {
		fmt.Println("node hash", node.Hash)
		tempNode := node
		for tempNode != nil {
			if tempNode.Parent == nil {
				// root node
				break
			} else if tempNode.RightNode != nil {
				node.ProofPathNodeHashes = append(node.ProofPathNodeHashes, tempNode.RightNode.Hash)
				fmt.Println("right", tempNode.RightNode.Hash)
			} else if tempNode.LeftNode != nil {
				node.ProofPathNodeHashes = append(node.ProofPathNodeHashes, tempNode.LeftNode.Hash)
				fmt.Println("left", tempNode.LeftNode.Hash)
			}
			fmt.Println("parent", tempNode.Parent.Hash)
			tempNode = tempNode.Parent
		}
	}

	fmt.Println("Proof Path", leafs[txIndex].ProofPathNodeHashes, len(leafs[txIndex].ProofPathNodeHashes), leafs[txIndex])

}

func buildTree(leafs []*Node) (*Node, []*Node) {
	groups := SplitNodeGroup(leafs)
	nodes := ComposeNode(groups)
	switch len(nodes) {
	case 1:
		return nodes[0], nodes // root node
	default:
		rootNode, parentNodes := buildTree(nodes)
		if len(nodes) == 2 {
			nodes[0].RightNode, nodes[1].LeftNode = nodes[1], nodes[0]
		}
		return rootNode, append(parentNodes, nodes...)
	}
}

func mockTxs(n int) []Transaction {
	txs := make([]Transaction, n)
	for i := 0; i < n; i++ {
		txs[i] = Transaction{txID: []byte(strconv.Itoa(i))}
	}
	return txs
}
