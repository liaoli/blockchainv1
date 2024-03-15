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

	//锁定的脚本
	ScriptPubKey string
	//接收的金额
	Value float64
}

type Transaction struct {
	//交易ID
	TXID []byte

	//交易输入，可能是多个
	TXInputs []TXInput

	//交易输入，可能是多个
	TXOutputs []TXOutput
}

func NewTransaction(from, to string, amount float64, bc BlockChain) *Transaction {
	//1。from/付款人/ to/收款人，amount/交易数量
	// 2。遍历账本，找到from满足条件的utxo集合，返回这些utxo包含的总金额
	//所有将要使用的utxo
	var spentUTXO = make(map[string][]int64)
	//所有将要使用的utxo的总额
	var retValue float64
	spentUTXO, retValue = bc.FindNeedUTXO(from, amount)
	fmt.Println(retValue)
	// 3。如果金额不足，创建交易失败
	if retValue < amount {
		fmt.Println("金额不足，创建交易失败")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput
	// 4。拼接inputs
	//>遍历utxo集合，每一个output都要转换为一个input
	for txid, indexArray := range spentUTXO {
		for _, i := range indexArray {
			input := TXInput{[]byte(txid), i, from}
			inputs = append(inputs, input)
		}
	}
	//5。拼接outputs
	//> 创建属于to 的outpt
	output1 := TXOutput{to, amount}

	outputs = append(outputs, output1)
	//>如果总额大于需要的转账金额，进行找零：给from创建output

	if retValue > amount {
		output2 := TXOutput{from, retValue - amount}

		outputs = append(outputs, output2)
	}
	// 6。设置哈希，返回

	//timeStamp := time.Now().Unix()

	tx := Transaction{
		nil,
		inputs,
		outputs,
		//timeStamp,
	}
	tx.setHash()
	return &tx
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

	output := TXOutput{address, reward}

	txTmp := Transaction{nil, []TXInput{input}, []TXOutput{output}}

	txTmp.setHash()

	return &txTmp

}

func (tx *Transaction) isCoinBaseTx() bool {

	inputs := tx.TXInputs
	//input个数为1，id为nil，索引为-1
	if len(inputs) == 1 && inputs[0].TXID == nil && inputs[0].VoutIndex == -1 {
		return true
	}
	return false
}
