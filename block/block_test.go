package block_test

import (
	"reflect"
	"testing"

	"github.com/cristaloleg/bc/block"
)

func TestGenesisBlock(t *testing.T) {
	gen := block.GenesisBlock()
	block := block.NewBlock(0, []byte("<put here your favourite words of wisdom>"), []byte{})

	if !reflect.DeepEqual(gen.Data, block.Data) {
		t.Fatalf("want Data %v, got %v", block.Data, gen.Data)
	}
	if !reflect.DeepEqual(gen.Hash, block.Hash) {
		t.Fatalf("want Hash %v, got %v", block.Hash, gen.Hash)
	}
	if !reflect.DeepEqual(gen.PrevHash, block.PrevHash) {
		t.Fatalf("want PrevHash %v, got %v", block.PrevHash, gen.PrevHash)
	}
}
