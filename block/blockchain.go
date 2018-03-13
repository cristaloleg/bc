package block

import (
	"log"
	"sync/atomic"
)

// Blockchain represents a blockchain
type Blockchain struct {
	storage     Storage
	proof       ProofOfWork
	lastHash    []byte
	lastBlockID int64
}

// Storage is an interface for storing blocks
type Storage interface {
	Init(*Block) ([]byte, error)
	GetBlock(hash []byte) (*Block, error)
	AddBlock(*Block) error
	GetLastHash() ([]byte, error)
	Close() error
}

// ProofOfWork is an interface for prooving work
type ProofOfWork interface {
	IsValid(b *Block) bool
	DoWork(b *Block) *Block
}

// NewBlockchain creates a new blockchain with a given params
func NewBlockchain(proof ProofOfWork, storage Storage) *Blockchain {
	bc := &Blockchain{
		proof:   proof,
		storage: storage,
	}

	genBlock := GenesisBlock()
	genBlock = bc.proof.DoWork(genBlock)
	bc.lastHash, _ = bc.storage.Init(genBlock)
	return bc
}

// Iterator returns an iterator for given blockchain
func (bc *Blockchain) Iterator() *Iterator {
	return NewIterator(bc, bc.lastHash)
}

// AddBlock creates a new block with a given data
func (bc *Blockchain) AddBlock(data []byte) (hash []byte, err error) {
	id := atomic.AddInt64(&bc.lastBlockID, 1)

	block := NewBlock(id, data, bc.lastHash)
	block = bc.proof.DoWork(block)

	bc.lastHash = block.Hash
	err = bc.storage.AddBlock(block)
	return block.Hash, err
}

// GetBlock returns a block with a given hash if possible
func (bc *Blockchain) GetBlock(hash []byte) (*Block, error) {
	block, err := bc.storage.GetBlock(hash)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return block, err
}

// Stop to a blockchain
func (bc *Blockchain) Stop() error {
	return bc.storage.Close()
}
