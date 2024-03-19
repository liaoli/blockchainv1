package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"math/big"
	"time"
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
	//ScriptSig string

	ScriptSig []byte //对应当前交易的签名

	PubKey []byte //付款公钥
}

//包含资金接收方的相关信息,包含：
//接收金额
//锁定脚本
//==易错点：经常把Value写成小写字母开头的==，这样会无法写入数据库，切记

type TXOutput struct {

	//锁定的脚本
	//ScriptPubKey string
	ScriptPubKeyHash []byte //收款人的公钥匙哈希
	//接收的金额
	Value float64
}

//NewTxOutput 由于没有办法直接将地址赋值给TXoutPut，所以需要提供一个output的方法
func NewTxOutput(address string, amount float64) TXOutput {
	output := TXOutput{Value: amount}
	pubKeyHush := getPubKeyHashFromAddress(address)

	output.ScriptPubKeyHash = pubKeyHush

	return output
}

type Transaction struct {
	//交易ID
	TXID []byte

	//交易输入，可能是多个
	TXInputs []TXInput

	//交易输入，可能是多个
	TXOutputs []TXOutput

	//时间戳
	TimeStamp uint64
}

func NewTransaction(from, to string, amount float64, bc BlockChain) *Transaction {
	//钱包就时在这里使用的，from => 钱包里面找到对应的wallet->私钥->签名

	wm := NewWalletManager()

	if wm == nil {
		fmt.Println("打开钱包失败")
		return nil
	}
	// 钱包里面找到对应的wallet
	wallet, ok := wm.Wallets[from]

	if !ok {
		fmt.Println("没有找到付款人地址")
		return nil
	}
	fmt.Println("找到付款人的私钥和公钥，准备创建交易")
	priKey := wallet.PriKey
	pubKey := wallet.PubKey
	//w我们所有的output 都是由公钥哈希锁定的，所以去查找付款人的output时，也需要提供付款人的公钥哈希
	pubKeyHash := getPubKeyHashFromPubKey(pubKey)

	//1。from/付款人/ to/收款人，amount/交易数量
	// 2。遍历账本，找到from满足条件的utxo集合，返回这些utxo包含的总金额
	//所有将要使用的utxo
	var spentUTXO = make(map[string][]int64)
	//所有将要使用的utxo的总额
	var retValue float64
	spentUTXO, retValue = bc.FindNeedUTXO(pubKeyHash, amount)
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
			input := TXInput{[]byte(txid), i, nil, pubKey}
			inputs = append(inputs, input)
		}
	}
	//5。拼接outputs
	//> 创建属于to 的outpt
	output1 := NewTxOutput(to, amount)

	outputs = append(outputs, output1)
	//>如果总额大于需要的转账金额，进行找零：给from创建output

	if retValue > amount {
		output2 := NewTxOutput(from, retValue-amount)

		outputs = append(outputs, output2)
	}
	// 6。设置哈希，返回

	//timeStamp := time.Now().Unix()

	tx := Transaction{
		nil,
		inputs,
		outputs,
		uint64(time.Now().Unix()),
	}
	tx.setHash()

	if !bc.signTransaction(&tx, priKey) {
		fmt.Println("签名失败")
		return nil
	}

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
	input := TXInput{nil, -1, nil, []byte(data)}

	output := NewTxOutput(address, reward)

	txTmp := Transaction{nil, []TXInput{input}, []TXOutput{output}, uint64(time.Now().Unix())}

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

//实现具体签名动作（copy，设置为空，签名动作）
//参数1：私钥
//参数2：inputs所引用的output所在的交易集合：
//>key:交易id
//>value:交易本身
func (tx *Transaction) sign(priKey *ecdsa.PrivateKey, preTxs map[string]*Transaction) bool {

	fmt.Println("具体交易签名")

	if tx.isCoinBaseTx() {
		fmt.Println("找到挖矿交易，无需签名")
		return true
	}
	//1.获取copy，PubKey，ScriptPubKey字段置空
	txCopy := tx.trimmedCopy()
	//2.遍历交易的inputs for
	for i, input := range txCopy.TXInputs {
		fmt.Printf("开始对input[%d]进行签名", i)

		preTx := preTxs[string(input.TXID)]
		if preTx == nil {
			return false
		}

		output := preTx.TXOutputs[input.VoutIndex]
		//>	获取引用的output的公钥哈希
		txCopy.TXInputs[i].PubKey = output.ScriptPubKeyHash
		//>对copy交易进行签名
		txCopy.setHash()

		hashData := txCopy.TXID //我们去签名的具体数据

		//开始签名

		r, s, err := ecdsa.Sign(rand.Reader, priKey, hashData)

		if err != nil {
			fmt.Println("签名错误：", err)
			return false
		}
		signature := append(r.Bytes(), s.Bytes()...)

		//>将数字签名赋值给原始的tx
		tx.TXInputs[i].ScriptSig = signature

		//>将input的pubKey字段设置为nil
		txCopy.TXInputs[i].PubKey = nil
	}
	fmt.Println("交易签名成功")
	//交易签名成
	return true
}

//trimmedCopy 修剪备用
func (tx *Transaction) trimmedCopy() *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	for _, input := range tx.TXInputs {
		in := TXInput{
			input.TXID,
			input.VoutIndex,
			nil,
			nil,
		}
		inputs = append(inputs, in)
	}

	outputs = tx.TXOutputs

	txCopy := Transaction{tx.TXID, inputs, outputs, tx.TimeStamp}

	return &txCopy
}

//具体校验逻辑
func (tx *Transaction) verify(prevTxs map[string]*Transaction) bool {
	//1. 获取交易副本txCopy
	txCopy := tx.trimmedCopy()
	//2. 遍历交易，inputs，
	for i, input := range tx.TXInputs {
		prevTx := prevTxs[string(input.TXID)]
		if prevTx == nil {
			return false
		}

		//3. 还原数据（得到引用output的公钥哈希）获取交易的哈希值
		output := prevTx.TXOutputs[input.VoutIndex]
		txCopy.TXInputs[i].PubKey = output.ScriptPubKeyHash
		txCopy.setHash()

		//清零环境, 设置为nil
		txCopy.TXInputs[i].PubKey = nil

		//具体还原的签名数据哈希值
		hashData := txCopy.TXID
		//签名
		signature := input.ScriptSig
		//公钥的字节流
		pubKey := input.PubKey

		//开始校验
		var r, s, x, y big.Int
		//r,s 从signature截取出来
		r.SetBytes(signature[:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])

		//x, y 从pubkey截取除来，还原为公钥本身
		x.SetBytes(pubKey[:len(pubKey)/2])
		y.SetBytes(pubKey[len(pubKey)/2:])
		curve := elliptic.P256()
		pubKeyRaw := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

		//进行校验
		res := ecdsa.Verify(&pubKeyRaw, hashData, &r, &s)
		if !res {
			fmt.Println("发现校验失败的input!")
			return false
		}
	}
	//4. 通过tx.ScriptSig, tx.PubKey进行校验
	fmt.Println("交易校验成功!")

	return true
}
