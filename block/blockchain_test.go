package block_test

import (
	"reflect"
	"testing"

	"github.com/cristaloleg/bc/block"
)

func TestBlockchain(t *testing.T) {
	proof, storage := setupTestBlockchain(t)
	bc := block.NewBlockchain(proof, storage)
	_ = bc

	genBlock := block.GenesisBlock()

	hash, err := bc.AddBlock([]byte("mydata"))
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if !reflect.DeepEqual(genBlock.Hash, hash) {
		t.Fatalf("want Hash %v, got %v", genBlock.Hash, hash)
	}
}

func setupTestBlockchain(t *testing.T) (block.ProofOfWork, block.Storage) {
	t.Helper()

	var proof = mockProofOfWork{
		IsValidFunc: func(b *block.Block) bool {
			return true
		},
		DoWorkFunc: func(b *block.Block) *block.Block {
			return b
		},
	}

	var storage = mockStorage{
		InitFunc: func(b *block.Block) ([]byte, error) {
			return []byte("three"), nil
		},
		GetLastHashFunc: func() ([]byte, error) {
			return []byte("three"), nil
		},
		GetBlockFunc: func(hash []byte) (*block.Block, error) {
			blocks := map[string]*block.Block{
				"one":   &block.Block{Hash: []byte("one"), PrevHash: nil},
				"two":   &block.Block{Hash: []byte("two"), PrevHash: []byte("one")},
				"three": &block.Block{Hash: []byte("three"), PrevHash: []byte("two")},
			}
			return blocks[string(hash)], nil
		},
		AddBlockFunc: func(b *block.Block) error {
			return nil
		},
	}

	return proof, storage
}
