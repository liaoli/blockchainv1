package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

type BlockChainIterator struct {
	db *bolt.DB
	//指向当前区块
	current_point []byte
}

func NewBlockChainIterator(bc *BlockChain) *BlockChainIterator {
	var it BlockChainIterator
	it.db = bc.db
	it.current_point = bc.tail

	return &it
}

func (it *BlockChainIterator) GetBlockAnMoveLeft() Block {
	var block Block

	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))

		if bucket == nil {
			fmt.Println("bucket should not be nil")
			os.Exit(1)
		} else {
			fmt.Printf("now : %x\n", it.current_point)
			current_block_temp := bucket.Get(it.current_point)
			current_block := DeSerialize(current_block_temp)

			block = current_block
			it.current_point = current_block.PreHash
			fmt.Printf("next : %x\n", it.current_point)
		}

		return nil
	})

	return block
}
