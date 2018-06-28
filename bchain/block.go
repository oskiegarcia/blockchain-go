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

	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			//fmt.Println(calculateHash(newBlock), " do more work")
			//time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(calculateHash(newBlock), fmt.Sprintf(" work done after %d iterations", i))
			newBlock.Hash = []byte(calculateHash(newBlock))
			break
		}
	}

	return newBlock
}

// GenesisBlock ...
func GenesisBlock() *Block {
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), 0, difficulty, fmt.Sprintf("%x", 100), []byte{}, []byte(calculateHash(&genesisBlock))}
	spew.Dump(genesisBlock)
	//Blockchain = append(Blockchain, genesisBlock)

	return &genesisBlock
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
