package block_test

import (
	"testing"

	"github.com/cristaloleg/bc/block"
)

func TestProof(t *testing.T) {
	proof := block.NewProof(3, 100)
	if proof == nil {
		t.Fail()
	}

	block := block.NewBlock(0, []byte("mydata"), nil)
	block = proof.DoWork(block)

	if !proof.IsValid(block) {
		t.Fail()
	}
}

var blockSink *block.Block

func BenchmarkProof(b *testing.B) {
	b.ReportAllocs()

	proof := block.NewProof(3, 100)
	for i := 0; i < b.N; i++ {
		blockSink = block.NewBlock(0, []byte("mydata"), nil)
		blockSink = proof.DoWork(blockSink)
	}
}

func BenchmarkProof5(b *testing.B) {
	b.ReportAllocs()

	proof := block.NewProof(5, 1000)
	for i := 0; i < b.N; i++ {
		blockSink = block.NewBlock(0, []byte("mydata"), nil)
		blockSink = proof.DoWork(blockSink)
	}
}
