package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

type BlockChain struct {
	//定义一个区块链的结构
	//Blocks []*Block
	//操作数据库的句柄
	db *bolt.DB
	//尾巴，存储最后一个区块的哈希
	tail []byte
}

//创世语
const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
const blockchainDBFile = "blockchain.db"
const bucketBlock = "bucketBlock"
const lastBlockHashKey = "lastBlockHashKey" //用于访问bolt数据库，得到最好一个区块的哈希值

func NewBlockChain() *BlockChain {
	//对区块链进行初始化，并把创世块添加到区块链
	var lastHash []byte
	//1. 打开数据库 没有的话创建
	//2.找到抽屉(bucket),如果找到，就返回bucket，如果没找到，就通过名字创建
	//    a.找到了
	//        通过'last'这个key找到我们做最后的一个哈希
	//    b.没找到
	//	     1.创建bucket，通过名字
	//         2.添加创世块数据
	//         3.更新"last"这个key的value(创世块的哈希值)

	db, err := bolt.Open(blockchainDBFile, 0600, nil)

	if err != nil {
		fmt.Println("bolt open failed", err)
		os.Exit(1)
	}

	db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(bucketBlock))
		var err error

		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(bucketBlock))
			if err != nil {
				fmt.Println("create  bucket failed", err)
				os.Exit(1)
			}
			genesisblock := NewBlock(genesisInfo, []byte{})
			bucket.Put(genesisblock.Hash, genesisblock.Serialize())
			bucket.Put([]byte(lastBlockHashKey), genesisblock.Hash)

			lastHash = genesisblock.Hash
		} else {
			lastHash = bucket.Get([]byte(lastBlockHashKey))
		}
		return nil
	})
	return &BlockChain{db, lastHash}
}

func (bc *BlockChain) AddBlock(data string) {

	//1.产生区块链
	//  数据和区块的哈希值
	//  通过数组的下标，拿到最后一个区块的哈希值，这个哈希值就是我们新区块的前哈希值
	//blockLen := len(bc.Blocks)
	//
	//lastBlock := bc.Blocks[blockLen-1]
	//prevBlockHash := lastBlock.Hash
	//block := NewBlock(data, prevBlockHash)
	//bc.Blocks = append(bc.Blocks, block)
	//获取最后一个hash区块
	lastBlockHash := bc.tail
	newBlock := NewBlock(data, lastBlockHash)

	bc.db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(bucketBlock))
		if bucket == nil {
			fmt.Println("bucket should not be nil!!")
			os.Exit(1)
		} else {
			bucket.Put(newBlock.Hash, newBlock.Serialize())
			bucket.Put([]byte(lastBlockHashKey), newBlock.Hash)
			bc.tail = newBlock.Hash
		}
		return nil
	})

}
