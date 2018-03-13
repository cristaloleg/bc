package block_test

import (
	"github.com/cristaloleg/bc/block"
)

type mockStorage struct {
	InitFunc        func(b *block.Block) ([]byte, error)
	GetBlockFunc    func(hash []byte) (*block.Block, error)
	AddBlockFunc    func(b *block.Block) error
	GetLastHashFunc func() ([]byte, error)
}

func (mock mockStorage) Init(b *block.Block) ([]byte, error) {
	return mock.InitFunc(b)
}

func (mock mockStorage) AddBlock(b *block.Block) error {
	return mock.AddBlockFunc(b)
}

func (mock mockStorage) GetBlock(hash []byte) (*block.Block, error) {
	return mock.GetBlockFunc(hash)
}

func (mock mockStorage) GetLastHash() ([]byte, error) {
	return mock.GetLastHashFunc()
}

func (mock mockStorage) Close() error {
	return nil
}

type mockProofOfWork struct {
	IsValidFunc func(b *block.Block) bool
	DoWorkFunc  func(b *block.Block) *block.Block
}

func (mock mockProofOfWork) IsValid(b *block.Block) bool {
	return mock.IsValidFunc(b)
}
func (mock mockProofOfWork) DoWork(b *block.Block) *block.Block {
	return mock.DoWorkFunc(b)
}
