package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/boltdb/bolt"
	"math/big"
	"os"
)

func main() {
	//fmt.Printf("hello world!")
	//block := GenesisBlock(genesisInfo, []byte{})
	//fmt.Printf("PreHash : %x\n", block.PreHash)
	//fmt.Printf("Hash : %x\n", block.Hash)
	//fmt.Printf("Data : %s\n", block.Data)

	//bc := NewBlockChain("test")
	cli := CLI{}
	cli.Run()
	//cli.CreateWallet()
	//ecdsaDemo()
	//cli.getBalance("liaoli")
	//cli.send("1DhvU59canfr4SwYLFuCA3qao7VLdNN4nE", "12EAzsXx9vWhLvhiNaA6kufwEc4csKEcdd", 4.5, "1DhvU59canfr4SwYLFuCA3qao7VLdNN4nE", "ok")
	//it := NewBlockChainIterator(bc)
	//
	//for {
	//	block := it.GetBlockAnMoveLeft()
	//	fmt.Println("============== ===============")
	//	fmt.Printf("PreHash : %x\n", block.PreHash)
	//	fmt.Printf("Hash : %x\n", block.Hash)
	//	fmt.Printf("Data : %s\n", block.Data)
	//	fmt.Printf("MerkleRoot : %x\n", block.MerkleRoot)
	//	fmt.Printf("Nonce : %d\n", block.Nonce)
	//	fmt.Printf("Version : %d\n", block.Version)
	//	//时间格式化
	//	timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
	//	//fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
	//	fmt.Printf("TimeStamp : %s\n", timeFormat)
	//	fmt.Printf("Difficulty : %d\n", block.Difficulty)
	//	pow := NewProofOfWork(block)
	//	fmt.Printf("IsValid:%v \n", pow.IsValid())
	//
	//	if len(block.PreHash) == 0 {
	//		break
	//	}
	//}
	//
	//for i, block := range bc.Blocks {
	//	fmt.Println("==============block height：", i, "===============")
	//	fmt.Printf("PreHash : %x\n", block.PreHash)
	//	fmt.Printf("Hash : %x\n", block.Hash)
	//	fmt.Printf("Data : %s\n", block.Data)
	//	fmt.Printf("MerkleRoot : %x\n", block.MerkleRoot)
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

}

func boltDemo() {
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

//go语言只提供了签名校验，未提供加解密
//创建私钥
//私钥得到公钥
//私钥签名
//公钥验证

func ecdsaDemo() {
	curve := elliptic.P256()
	//创建私钥
	priKey, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		fmt.Println("ecdsa.GenerateKey err", err)
	}

	//fmt.Println("PriKey=", PriKey)
	//私钥得到公钥
	pubKey := priKey.PublicKey
	data := "hello world"
	hash := sha256.Sum256([]byte(data))
	//私钥签名
	r, s, err := ecdsa.Sign(rand.Reader, priKey, hash[:])
	if err != nil {
		fmt.Println("ecdsa.Sign err", err)
		return
	}

	signature := append(r.Bytes(), s.Bytes()...)
	//传输。。。
	//

	//公钥验证
	//在对端将r,s取出来
	var r1, s1 big.Int
	r1.SetBytes(signature[:len(signature)/2])
	s1.SetBytes(signature[len(signature)/2:])
	reuslt := ecdsa.Verify(&pubKey, hash[:], &r1, &s1)

	fmt.Println("ecdsa.Verify result=", reuslt)
}
