package block

import (
	"log"
)

// Iterator represents an interator over a Blockchain
type Iterator struct {
	bc       *Blockchain
	currHash []byte
	isDone   bool
}

// NewIterator returns a new iterator for a Blockchain pointing on given hash
func NewIterator(bc *Blockchain, hash []byte) *Iterator {
	it := &Iterator{
		bc:       bc,
		currHash: hash,
		isDone:   false,
	}
	return it
}

// Next if possible gp to next block and return it, nil otherwise
func (it *Iterator) Next() *Block {
	if it.isDone {
		return nil
	}
	block, err := it.bc.GetBlock(it.currHash)

	if err != nil || block == nil {
		log.Println(err)
		it.isDone = true
		return nil
	}

	it.currHash = block.PrevHash
	return block
}

// HasNext returns false when you cannot iterate further
func (it *Iterator) HasNext() bool {
	return !it.isDone
}
