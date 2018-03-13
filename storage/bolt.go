package storage

import (
	"encoding/json"
	"log"
	"time"

	"github.com/cristaloleg/bc/block"

	"github.com/coreos/bbolt"
)

var _ block.Storage = (*Bolt)(nil)

// Bolt represents a BoltDB
type Bolt struct {
	client *bolt.DB
	table  string
}

// NewBolt will instantiate a BoltDB client
func NewBolt(filename string) *Bolt {
	c, err := bolt.Open(filename, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		panic(err)
	}

	db := &Bolt{
		client: c,
		table:  "blocks",
	}
	return db
}

// Init will initialize a database
func (db *Bolt) Init(block *block.Block) (lastHash []byte, err error) {
	err = db.client.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.table))

		if b != nil {
			lastHash = b.Get([]byte("lastHash"))
			return nil
		}

		lastHash = block.Hash

		b, err := tx.CreateBucket([]byte(db.table))
		if err != nil {
			log.Panic(err)
			return err
		}

		raw, err := json.Marshal(block)
		if err != nil {
			return nil
		}
		err = b.Put(block.Hash, raw)
		if err != nil {
			log.Panic(err)
			return err
		}

		err = b.Put([]byte("lastHash"), block.Hash)
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	return lastHash, err
}

// GetBlock returns a block by it's hash
func (db *Bolt) GetBlock(hash []byte) (bloc *block.Block, err error) {
	err = db.client.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.table))
		raw := b.Get(hash)
		var bl block.Block
		err := json.Unmarshal(raw, &bl)
		bloc = &bl
		return err
	})
	return bloc, err
}

// GetLastHash returns a latest hash in a blockchain
func (db *Bolt) GetLastHash() ([]byte, error) {
	var lastHash []byte
	err := db.client.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.table))
		lastHash = b.Get([]byte("lastHash"))
		return nil
	})
	return lastHash, err
}

// AddBlock will save block in the database
func (db *Bolt) AddBlock(block *block.Block) error {
	raw, err := json.Marshal(block)
	if err != nil {
		log.Panic(err)
	}

	err = db.client.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.table))

		if err := b.Put(block.Hash, raw); err != nil {
			log.Panic(err)
		}

		if err := b.Put([]byte("lastHash"), block.Hash); err != nil {
			log.Panic(err)
		}
		return nil
	})
	return err
}

// Close save data and close the connection
func (db *Bolt) Close() error {
	return db.client.Close()
}
