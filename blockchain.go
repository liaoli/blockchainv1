package main

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
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

func NewBlockChain(address string) *BlockChain {

	//1.区块拉链不存在，创建
	if isFileExist(blockchainDBFile) {
		fmt.Println("区块链文件已经存在")
		return nil
	}
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
			//创建coinbase交易
			coinbaseTx := NewCoinbaseTx(address, genesisInfo)
			genesisBlock := NewBlock([]*Transaction{coinbaseTx}, []byte{})
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)

			lastHash = genesisBlock.Hash
		} else {
			lastHash = bucket.Get([]byte(lastBlockHashKey))
		}
		return nil
	})
	return &BlockChain{db, lastHash}
}

func GetBlockChainInstance() (*BlockChain, error) {
	if !isFileExist(blockchainDBFile) {
		fmt.Println("区块链文件不存在，请先创建")
		return nil, errors.New("区块链文件不存在，请先创建")
	}

	var lastHash []byte //内存中最后一个区块的哈希值

	//两个功能
	//如果区块链不存在，则创建，同时返回blockchain的示例
	db, err := bolt.Open(blockchainDBFile, 0600, nil)

	if err != nil {
		return nil, err
	}

	//区块链存在则直接返回blockchain对象
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))
		if bucket == nil {
			return errors.New("bucket 不存在")
		} else {
			lastHash = bucket.Get([]byte(lastBlockHashKey))
		}

		return nil
	})

	bc := BlockChain{db: db, tail: lastHash}
	return &bc, nil
}
func (bc *BlockChain) AddBlock(tx []*Transaction) error {

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
	newBlock := NewBlock(tx, lastBlockHash)

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

	return nil
}

func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return NewBlockChainIterator(bc)
}

//UTXOInfo 包含output本身，位置信息
type UTXOInfo struct {
	//交易id
	TxId []byte
	//索引值
	index int64
	//output
	TXOutput
}

//获取指定地址的金额,实现遍历账本的通用函数
//给定一个地址，返回所有的utxo

//FindMyUTXO 返回制定地址能够支配的utxo所在的交易集合
func (bc *BlockChain) FindMyUTXO(pubKeyHash []byte) []UTXOInfo {
	//var transactions []Transaction
	//var outputs []UTXOInfo
	var utxoInfos []UTXOInfo
	spentUTXOs := make(map[string][]int64)
	it := bc.NewIterator()
	for {
		//1.遍历区块
		block := it.GetBlockAnMoveLeft()
		//2.遍历交易
		for _, tx := range block.Transactions {
			//3.遍历outputs 判断这个output的锁定脚本是否为我们的目标地址
		LABEL:
			for outputIndex, output := range tx.TXOutputs {

				//if output.ScriptPubKeyHash == pubKeyHash {
				if bytes.Equal(output.ScriptPubKeyHash, pubKeyHash) {
					fmt.Println("outputIndex:", outputIndex)
					//开始过滤
					currentTxid := string(tx.TXID)

					indexArray := spentUTXOs[currentTxid]

					if len(indexArray) != 0 {
						for _, spendIndex := range indexArray {
							if outputIndex == int(spendIndex) {
								continue LABEL
							}
						}
					}
					//outputs = append(outputs, output)
					utxoInfo := UTXOInfo{tx.TXID, int64(outputIndex), output}

					utxoInfos = append(utxoInfos, utxoInfo)
				}

			}
			if !tx.isCoinBaseTx() {
				//非挖矿交易才遍历
				for inputIndex, input := range tx.TXInputs {
					if bytes.Equal(getPubKeyHashFromPubKey(input.PubKey), pubKeyHash) /*付款人的公钥*/ {
						fmt.Println("inputIndex:", inputIndex)
						spentKey := string(input.TXID)
						spentUTXOs[spentKey] = append(spentUTXOs[spentKey], input.VoutIndex)
					}

				}
			}

		}

		if len(block.PreHash) == 0 {
			break
		}
	}

	return utxoInfos
}

func (bc *BlockChain) FindNeedUTXO(pubKeyHash []byte, amount float64) (map[string][]int64, float64) {

	//两个返回值
	retMap := make(map[string][]int64)
	var retValue float64
	//1.遍历账本，找到所有utxo
	utxoInfos := bc.FindMyUTXO(pubKeyHash)
	//2.遍历utxo，统计当前总额，与amount比较
	for _, utxoInfo := range utxoInfos {
		//统计当前utxo的总和
		retValue += utxoInfo.Value
		key := string(utxoInfo.TxId)
		retMap[key] = append(retMap[key], utxoInfo.index)
		if retValue >= amount {
			break

		}
		//>如果大于等于amount直接返回
		//>反之继续遍历
	}

	return retMap, retValue
}

//交易签名函数
func (bc BlockChain) signTransaction(tx *Transaction, priKey *ecdsa.PrivateKey) bool {
	fmt.Println("开始签名交易")
	//根据传递进来的tx 得到所有需要的前交易preTxs
	preTxs := make(map[string]*Transaction)
	//遍历账本，找到所有需要的交易集合
	for _, input := range tx.TXInputs {
		pretx := bc.findTransaction(input.TXID)

		if pretx == nil {
			fmt.Println("没有找到有效引用的交易")
			return false
		}

		fmt.Println("到有效引用的交易")
		//容易错误：tx.TXID
		preTxs[string(input.TXID)] = pretx
	}

	return tx.sign(priKey, preTxs)
}

func (bc *BlockChain) findTransaction(txid []byte) *Transaction {
	//遍历区块，遍历账本，比较txid与交易id，如果相同，返回交易，反之返回nil
	it := bc.NewIterator()

	for {
		block := it.GetBlockAnMoveLeft()
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.TXID, txid) {
				return tx
			}
		}

		if len(block.PreHash) == 0 {
			break
		}
	}
	return nil
}
