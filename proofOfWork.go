package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//- 定⼀个⼯作量证明的结构ProofOfWork
//- block
//- ⽬标值
type ProofOfWork struct {
	//区块数据
	block *Block
	//⽬标值，先写成固定的值，后⾯再进⾏推到演算。
	target big.Int
}

// 创建工作量证明的函数
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}
	//⾃定义的难度值，先写成固定值
	//0000100000000000000000000000000000000000000000000000000000000000
	targetString := "0000100000000000000000000000000000000000000000000000000000000000"

	bigIntTmp := big.Int{}
	bigIntTmp.SetString(targetString, 16)

	pow.target = bigIntTmp

	return &pow
}

func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var Nonce uint64
	var hash [32]byte
	fmt.Printf("target : %s \n", pow.target.String())

	for {
		hash = sha256.Sum256(pow.prepareDate(Nonce))

		tTmp := big.Int{}
		tTmp.SetBytes(hash[:])
		//当前哈希.Cmp(⽬标值) == -1，说明当前的哈希⼩于⽬标值
		if tTmp.Cmp(&pow.target) == -1 {
			fmt.Printf("found hash :%x ,%d \n", hash, Nonce)
			break
		} else {
			Nonce++
		}
	}

	return hash[:], Nonce
}

func (pow *ProofOfWork) prepareDate(num uint64) []byte {
	block := pow.block
	//block.HashTransactionMerkleRoot()
	temp := [][]byte{
		block.PreHash,
		//block.Data,
		block.MerkleRoot,
		uint64ToByte(block.Version),
		uint64ToByte(block.TimeStamp),
		uint64ToByte(block.Difficulty),
		uint64ToByte(num),
	}
	data := bytes.Join(temp, []byte(""))

	return data
}

func (pow *ProofOfWork) IsValid() bool {
	hash := sha256.Sum256(pow.prepareDate(pow.block.Nonce))
	fmt.Printf("is valid hash : %x ,%d \n", hash[:], pow.block.Nonce)

	tTmp := big.Int{}
	tTmp.SetBytes(hash[:])
	if tTmp.Cmp(&pow.target) == -1 {
		return true
	}
	return false
}
