package storage

import (
	"encoding/json"
	"log"

	"github.com/cristaloleg/bc/block"

	"github.com/tidwall/buntdb"
)

var _ block.Storage = (*BuntDB)(nil)

// BuntDB represents a BuntDB
type BuntDB struct {
	client *buntdb.DB
}

// NewBuntDB will instantiate a BoltDB client
func NewBuntDB(filename string) *BuntDB {
	c, err := buntdb.Open(filename)
	if err != nil {
		panic(err)
	}

	db := &BuntDB{
		client: c,
	}
	return db
}

// Init will initialize a database
func (db *BuntDB) Init(block *block.Block) (lastHash []byte, err error) {
	err = db.client.Update(func(tx *buntdb.Tx) error {
		savedHash, err := tx.Get("lastHash")
		if err == nil {
			lastHash = []byte(savedHash)
			return nil
		}
		lastHash = []byte(block.Hash)

		raw, err := json.Marshal(block)
		if err != nil {
			return nil
		}
		_, _, err = tx.Set(string(block.Hash), string(raw), nil)
		if err != nil {
			log.Panic(err)
			return err
		}

		_, _, err = tx.Set(("lastHash"), string(block.Hash), nil)
		if err != nil {
			return err
		}

		return nil
	})
	return lastHash, err
}

// GetBlock returns a block by it's hash
func (db *BuntDB) GetBlock(hash []byte) (block *block.Block, err error) {
	err = db.client.View(func(tx *buntdb.Tx) error {
		raw, err := tx.Get(string(hash))
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(raw), block)
		return err
	})
	return block, err
}

// GetLastHash returns a latest hash in a blockchain
func (db *BuntDB) GetLastHash() ([]byte, error) {
	var lastHash []byte
	err := db.client.View(func(tx *buntdb.Tx) error {
		hash, err := tx.Get("lastHash")
		lastHash = []byte(hash)
		return err
	})
	return lastHash, err
}

// AddBlock will save block in the database
func (db *BuntDB) AddBlock(block *block.Block) error {
	err := db.client.Update(func(tx *buntdb.Tx) error {
		var err error

		raw, err := json.Marshal(block)
		if err != nil {
			return nil
		}

		_, _, err = tx.Set(string(block.Hash), string(raw), nil)
		if err != nil {
			return err
		}

		_, _, err = tx.Set(("lastHash"), string(block.Hash), nil)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// Close save data and close the connection
func (db *BuntDB) Close() error {
	return db.client.Close()
}
