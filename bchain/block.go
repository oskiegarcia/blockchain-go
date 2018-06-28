package bchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index      int
	Timestamp  string
	BPM        int
	Difficulty int
	Nonce      string
	PrevHash   []byte
	Hash       []byte
}

// GenerateBlock ...
func (bc *Blockchain) GenerateBlock(BPM int) *Block {

	log.Printf("GenerateBlock: %d", BPM)

	newBlock := &Block{}

	lastBlock := bc.lastBlock()
	log.Printf("Last Index: %d", lastBlock.Index)

	t := time.Now()
	newBlock.Index = lastBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = lastBlock.Hash
	newBlock.Difficulty = difficulty
	newBlock.proofOfWork()

	return newBlock
}

// GenesisBlock - creates the initial block
func GenesisBlock() *Block {
	t := time.Now()
	genesisBlock := &Block{0, t.String(), 0, difficulty, fmt.Sprintf("%x", 100), []byte{}, []byte{}}
	genesisBlock.proofOfWork()
	spew.Dump(genesisBlock)

	return genesisBlock
}

// proofOfWork
func (b *Block) proofOfWork() {

	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		b.Nonce = hex
		if !isHashValid(calculateHash(b), b.Difficulty) {
			continue
		} else {
			fmt.Println(calculateHash(b), fmt.Sprintf(" work done after %d iterations", i))
			b.Hash = []byte(calculateHash(b))
			break
		}
	}
}

// Serialize serializes the block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock deserializes a block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
