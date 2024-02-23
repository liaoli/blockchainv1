package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

//引用utxo所在交易的ID
//所消费utxo在output中的索引
//解锁脚本

type TXInput struct {
	//引用output所在交易ID
	TXID []byte

	//引用output的索引值
	VoutIndex int64

	//解锁脚本
	ScriptSig string
}

//包含资金接收方的相关信息,包含：
//接收金额
//锁定脚本
//==易错点：经常把Value写成小写字母开头的==，这样会无法写入数据库，切记

type TXOutput struct {
	//接收的金额
	Value float64

	//锁定的脚本
	ScriptPubKey string
}

type Transaction struct {
	//交易ID
	TXID []byte

	//交易输入，可能是多个
	TXInputs []TXInput

	//交易输入，可能是多个
	TXOutputs []TXOutput
}

//设置交易ID方法
func (t Transaction) setHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(t)
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	t.TXID = hash[:]
}

//挖矿奖励
const reward = 12.5

// NewCoinbaseTx 创建Coinbase交易
func NewCoinbaseTx(address string, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("reward %s %f", address, reward)
	}
	////比特币系统，对于这个input的id填0，对索引填0xffff，data由矿工填写，一般填所在矿池的名字
	input := TXInput{nil, -1, data}

	output := TXOutput{reward, address}

	txTmp := Transaction{nil, []TXInput{input}, []TXOutput{output}}

	txTmp.setHash()

	return &txTmp

}
