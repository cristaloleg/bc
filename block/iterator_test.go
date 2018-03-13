package block_test

import (
	"testing"

	"github.com/cristaloleg/bc/block"
)

func TestIterator(t *testing.T) {
	proof, storage := setupTestIterator(t)
	bc := block.NewBlockchain(proof, storage)

	it := bc.Iterator()
	if it == nil {
		t.Fail()
	}

	var b *block.Block
	hashes := []string{"three", "two", "one"}

	for _, hash := range hashes {
		b = it.Next()
		if string(b.Hash) != hash {
			t.Fatalf("want %v, got %v", hash, string(b.Hash))
		}
	}

	b = it.Next()
	if b != nil {
		t.Fatal("expected to be nil", string(b.Hash))
	}
	b = it.Next()
	if b != nil {
		t.Fatal("expected to be nil", string(b.Hash))
	}
}

func setupTestIterator(t *testing.T) (block.ProofOfWork, block.Storage) {
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
	}

	return proof, storage
}
