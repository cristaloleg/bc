package block

import "time"

// Block represents a block in a blockchain
type Block struct {
	ID        int64  `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Data      []byte `json:"data"`
	Hash      []byte `json:"hash"`
	PrevHash  []byte `json:"prev_hash"`
	Nonce     int
}

// GenesisBlock creates a `genesis block`
func GenesisBlock() *Block {
	data := []byte("<put here your favourite words of wisdom>")
	return NewBlock(0, data, []byte{})
}

// NewBlock creates a new block in blockchain
func NewBlock(id int64, data, prevBlockHash []byte) *Block {
	b := &Block{
		ID:        id,
		Timestamp: time.Now().Unix(),
		Data:      data,
		Hash:      []byte{},
		PrevHash:  prevBlockHash,
		Nonce:     0,
	}
	return b
}
