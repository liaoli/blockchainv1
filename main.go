package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

func main() {
	//fmt.Printf("hello world!")
	//block := GenesisBlock(genesisInfo, []byte{})
	//fmt.Printf("PreHash : %x\n", block.PreHash)
	//fmt.Printf("Hash : %x\n", block.Hash)
	//fmt.Printf("Data : %s\n", block.Data)

	//bc := NewBlockChain()
	//bc.AddBlock("aa send 1 btc to c")
	//bc.AddBlock("bb send 1 btc to c")
	//
	//for i, block := range bc.Blocks {
	//	fmt.Println("==============block height：", i, "===============")
	//	fmt.Printf("PreHash : %x\n", block.PreHash)
	//	fmt.Printf("Hash : %x\n", block.Hash)
	//	fmt.Printf("Data : %s\n", block.Data)
	//	fmt.Printf("merkleRoot : %x\n", block.merkleRoot)
	//	fmt.Printf("Nonce : %d\n", block.Nonce)
	//	fmt.Printf("Version : %d\n", block.Version)
	//	//时间格式化
	//	timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
	//	//fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
	//	fmt.Printf("TimeStamp : %s\n", timeFormat)
	//	fmt.Printf("Difficulty : %d\n", block.Difficulty)
	//	pow := NewProofOfWork(*block)
	//	fmt.Printf("IsValid:%v \n", pow.IsValid())
	//}

	//1. 打开数据库
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		fmt.Println("bolt.Open failed!", err)
		os.Exit(1)
	}
	//2.写数据库
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("firstBucket"))
		var err error
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte("firstBucket"))
			if err != nil {
				fmt.Println("createBucket failed!", err)
				os.Exit(1)
			}
		}
		bucket.Put([]byte("aaaa"), []byte("HelloWorld!"))
		bucket.Put([]byte("bbbb"), []byte("HelloItcast!"))
		return nil
	})
	//3.读取数据库
	var value []byte
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("firstBucket"))
		if bucket == nil {
			fmt.Println("Bucket is nil!")
			os.Exit(1)
		}
		value = bucket.Get([]byte("aaaa"))
		fmt.Println("aaaa => ", string(value))
		value = bucket.Get([]byte("bbbb"))
		fmt.Println("bbbb => ", string(value))
		return nil
	})
}
