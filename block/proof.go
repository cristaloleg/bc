package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math/big"
)

var _ ProofOfWork = (*Proof)(nil)

// Proof is a Hashcash proof-of-work algorithm
type Proof struct {
	target     *big.Int
	targetBits int
	maxNonce   int
}

// NewProof instantiate new Hashcash proow-of-work
func NewProof(targetBits, maxNonce int) *Proof {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	p := &Proof{
		target:     target,
		targetBits: targetBits,
		maxNonce:   maxNonce,
	}
	return p
}

// IsValid returns true if block is valid, false otherwise
func (p *Proof) IsValid(b *Block) bool {
	var hashInt big.Int
	hashInt.SetBytes(p.getHash(b))
	return hashInt.Cmp(p.target) == -1
}

// DoWork process block until work is done
func (p *Proof) DoWork(b *Block) *Block {
	var hashInt big.Int

	for nonce := 0; nonce <= p.maxNonce; {
		b.Nonce = nonce
		hashInt.SetBytes(p.getHash(b))

		if hashInt.Cmp(p.target) == -1 {
			b.Hash = p.getHash(b)
			break
		}
		nonce++
	}
	return b
}

func (p *Proof) prepareData(b *Block, nonce int) []byte {
	return bytes.Join(
		[][]byte{
			intToHex(int64(p.targetBits)),
			intToHex(int64(nonce)),
			intToHex(b.Timestamp),
			b.Data,
			b.PrevHash,
		},
		[]byte{},
	)
}

func (p *Proof) getHash(b *Block) []byte {
	data := p.prepareData(b, b.Nonce)
	hash := sha256.Sum256(data)
	return hash[:]
}

func intToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes()
}
