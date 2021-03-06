package bchain

import (
	"errors"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

// Blockchain keeps a sequence of Blocks
type Blockchain struct {
	tip []byte
	Db  *bolt.DB
}

// BlockchainIterator is used to iterate over blockchain blocks
type BlockchainIterator struct {
	currentHash []byte
	Db          *bolt.DB
}

// BlockCollection ...
type BlockCollection struct {
	Size   int
	Blocks []*Block
}

// difficulty ...
const difficulty = 4

//-----------------------------------------
// Exported functions
//-----------------------------------------

// AddBlock - saves the block in the blockchain database
func (bc *Blockchain) AddBlock(newBlock *Block) error {

	log.Printf("AddBlock: %d\n", newBlock.Index)

	if !bc.isBlockValid(newBlock) {
		return errors.New("Invalid block")
	}

	err := bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})

	return err
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain() *Blockchain {

	log.Println("NewBlockchain...")

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := GenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash, genesis.serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}
	log.Println(string(tip))
	return &bc
}

// List - retrieves all data from the blockchain database
func (bc *Blockchain) List() *BlockCollection {
	bci := bc.iterator()
	var blockSlice []*Block
	result := BlockCollection{}
	for {
		block := bci.next()
		blockSlice = append(blockSlice, block)
		//log.Printf("hash:%s\n", block.Hash)
		if len(block.PrevHash) == 0 {
			break
		}
	}

	result.Blocks = blockSlice
	result.Size = len(blockSlice)
	return &result
}

//-----------------------------------------
// Unexported functions
//-----------------------------------------

// iterator - creates an iterator
func (bc *Blockchain) iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.Db}

	return bci
}

// next returns next block starting from the tip
func (i *BlockchainIterator) next() *Block {
	var block *Block

	err := i.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = deserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevHash

	return block
}

// returns the last block
func (bc *Blockchain) lastBlock() *Block {

	var block *Block

	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte("l"))
		log.Printf("lastHash:%s\n", string(lastHash))
		encodedBlock := b.Get(lastHash)
		block = deserializeBlock(encodedBlock)
		//spew.Dump(block)
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return block

}

// checks if it is a valid block
func (bc *Blockchain) isBlockValid(newBlock *Block) bool {

	lastBlock := bc.lastBlock()

	if lastBlock.Index+1 != newBlock.Index {
		return false
	}

	if string(lastBlock.Hash) != string(newBlock.PrevHash) {
		return false
	}

	if calculateHash(newBlock) != string(newBlock.Hash) {
		return false
	}

	return true
}
