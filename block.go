package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

type Block struct {
	//版本号
	Version uint64
	//前区块哈希值
	PreHash []byte
	//梅克尔根
	merkleRoot []byte
	//时间戳
	TimeStamp uint64
	//难度值
	Difficulty uint64
	//随机数，这就是挖矿时要寻求的数
	Nonce uint64
	//当前区块哈希值(为了⽅便实现，所以将区块的哈希值放到了区块中)	Hash []byte
	Hash []byte
	//区块数据
	Data []byte
}

func NewBlock(data string, preBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PreHash:    preBlockHash,
		merkleRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 100,
		Nonce:      100,
		Hash:       []byte{},
		Data:       []byte(data),
	}

	//block.SetHash()
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

func (block *Block) Serialize() []byte {
	//1.将block数据转换成流字节
	var buffer bytes.Buffer
	//创建一个编码器
	encoder := gob.NewEncoder(&buffer)
	//编码，将block编码城buffer
	err := encoder.Encode(block)

	if err != nil {
		fmt.Println("encode fialed", err)
		os.Exit(1)
	}
	return buffer.Bytes()

}

func DeSerialize(data []byte) Block {
	var block Block
	var buffer bytes.Buffer
	_, err := buffer.Write(data)
	if err != nil {
		fmt.Println("buffer.Write fialed", err)
		os.Exit(1)
	}
	decoder := gob.NewDecoder(&buffer)

	err = decoder.Decode(&block)

	if err != nil {
		fmt.Println("decode fialed", err)
		os.Exit(1)
	}

	return block
}

//产⽣创世块
func GenesisBlock(data string, prvBlockHash []byte) *Block {
	return NewBlock(data, prvBlockHash)
}

func (block *Block) SetHash() {
	var blockByteInfo []byte
	//1.拼接当前区块的数据
	//blockByteInfo = append(blockByteInfo, block.PreHash...)
	//blockByteInfo = append(blockByteInfo, block.Data...)
	//blockByteInfo = append(blockByteInfo, block.merkleRoot...)
	//blockByteInfo = append(blockByteInfo, uint64ToByte(block.Version)...)
	//blockByteInfo = append(blockByteInfo, uint64ToByte(block.TimeStamp)...)
	//blockByteInfo = append(blockByteInfo, uint64ToByte(block.Difficulty)...)
	//blockByteInfo = append(blockByteInfo, uint64ToByte(block.Nonce)...)

	//用Join 替代
	temp := [][]byte{
		block.PreHash,
		block.Data,
		block.merkleRoot,
		uint64ToByte(block.Version),
		uint64ToByte(block.TimeStamp),
		uint64ToByte(block.Difficulty),
		uint64ToByte(block.Nonce),
	}
	blockByteInfo = bytes.Join(temp, []byte(""))
	//2.对数据进行hash处理
	hash := sha256.Sum256(blockByteInfo)
	//3.把hash 添加到 Hash字段
	block.Hash = hash[:]

}
