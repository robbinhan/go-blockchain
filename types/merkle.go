package types

import (
	"crypto/sha256"
	"encoding/json"
)

type Node struct {
	Hash []byte // node hash
	// 两个子节点的hash
	ChildLeftNode       *Node // 左侧子节点
	ChildRightNode      *Node // 右侧子节点
	LeftNode            *Node // 左侧兄弟节点
	RightNode           *Node // 右侧兄弟节点
	Parent              *Node
	Data                []byte   // leaf 存储的数据
	ProofPathNodeHashes [][]byte // 当前节点到root节点的路径hash
}
type NodeGroup struct {
	Nodes []*Node
}

// SplitNodeGroup 将block的交易或者节点两两分组，最后不够的copy补位
func SplitNodeGroup(txs []*Node) []*NodeGroup {
	var groups []*NodeGroup

	for {
		var nodes []*Node
		if len(txs) >= 2 {
			nodes = txs[:2]
			txs = txs[2:]

			nodes[0].RightNode, nodes[1].LeftNode = nodes[1], nodes[0]
			groups = append(groups, &NodeGroup{
				Nodes: nodes,
			})
		} else {
			if len(txs) > 0 {
				nodes = []*Node{txs[0], txs[0]}
				nodes[0].RightNode, nodes[1].LeftNode = nodes[1], nodes[0]
				groups = append(groups, &NodeGroup{
					Nodes: nodes,
				})
			}
			break
		}
	}
	return groups
}

// ComposeNode 将每组node组合成新node
func ComposeNode(groups []*NodeGroup) []*Node {
	var nodes []*Node

	for _, group := range groups {
		newNodeHash := hashTwoNode(group.Nodes[0].Hash, group.Nodes[1].Hash)
		newNode := NewNode(nil)
		newNode.Hash = newNodeHash
		newNode.ChildLeftNode = group.Nodes[0]
		group.Nodes[0].Parent = newNode
		newNode.ChildRightNode = group.Nodes[1]
		group.Nodes[1].Parent = newNode
		nodes = append(nodes, newNode)
	}

	return nodes
}

// TreeNodes 递归将节点组合
func rootHash(leafs []*Node) []byte {
	groups := SplitNodeGroup(leafs)
	nodes := ComposeNode(groups)
	if len(nodes) > 1 {
		return rootHash(nodes)
	} else {
		return nodes[0].Hash
	}
}

func RootHash(block Block) []byte {
	var leafs []*Node

	for _, tx := range block.Txs {
		l := NewNode(tx.txID)
		leafs = append(leafs, l)
	}
	return rootHash(leafs)
}

func NewNode(data []byte) *Node {
	node := &Node{
		Data: data,
	}
	node.Hash = hash(node)
	return node
}

func hashTwoNode(nHash1, nHash2 []byte) []byte {
	shaHasher := sha256.New()

	hm := map[string]interface{}{
		"n1": nHash1,
		"n2": nHash2,
	}
	hmBytes, _ := json.Marshal(hm)

	shaHasher.Write(hmBytes)
	return shaHasher.Sum(nil)
}

func hash(n *Node) []byte {
	shaHasher := sha256.New()
	hmBytes, _ := json.Marshal(n.Data)

	shaHasher.Write(hmBytes)
	return shaHasher.Sum(nil)
}
