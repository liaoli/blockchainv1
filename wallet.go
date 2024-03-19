package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

//-结构定义
type Wallet struct {
	PriKey *ecdsa.PrivateKey

	PubKey []byte // xy 拼接而成 r,s
}

//-创建密钥对

func NewWalletKeyPair() *Wallet {
	//创建私钥
	curve := elliptic.P256()
	//创建私钥
	priKey, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		fmt.Println("ecdsa.GenerateKey err", err)
		return nil
	}
	//拼接公钥
	pubKeyRaw := priKey.PublicKey

	pubKey := append(pubKeyRaw.X.Bytes(), pubKeyRaw.Y.Bytes()...)
	//创建wallet钱包

	wallet := Wallet{PriKey: priKey, PubKey: pubKey}

	return &wallet
}

//-根据私钥生成地址
func (w *Wallet) getAddress() string {
	//公钥
	pubKey := w.PubKey
	//hash1 := sha256.Sum256(PubKey)
	////hash160处理
	////hasher := crypto.RIPEMD160.New()
	//
	//hasher := ripemd160.New()
	//
	//hasher.Write(hash1[:])
	////公钥哈希，锁定output时就是用这个值
	//pubKeyHash := hasher.Sum(nil)
	//拼接version和公钥哈希，得到21个字节

	pubKeyHash := getPubKeyHashFromPubKey(pubKey)

	payload := append([]byte{byte(00)}, pubKeyHash...)

	//生成4字节的校验码
	//first := sha256.Sum256(payload)
	//second := sha256.Sum256(first[:])
	//4字节checksum
	checkSum := checkSum(payload)

	payload = append(payload, checkSum...)

	address := base58.Encode(payload)

	return address
}

//给定公钥匙获取公钥的hash
func getPubKeyHashFromPubKey(pubKey []byte) []byte {
	hash1 := sha256.Sum256(pubKey)
	//hash160处理
	//hasher := crypto.RIPEMD160.New()

	hasher := ripemd160.New()

	hasher.Write(hash1[:])
	//公钥哈希，锁定output时就是用这个值
	pubKeyHash := hasher.Sum(nil)

	return pubKeyHash
}

//得到4字节的checkSum
func checkSum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	//4字节checksum
	checkSum := second[0:4]

	return checkSum

}

//通过地址，反推出公钥哈希，注意不是公钥
func getPubKeyHashFromAddress(address string) []byte {
	//base58解码

	decodeInfo := base58.Decode(address)
	//校验一下地址
	if len(decodeInfo) != 25 {
		fmt.Println("传入地址无效")
		return nil
	}

	//截取
	pubKeyHash := decodeInfo[1 : len(decodeInfo)-4]

	return pubKeyHash
}

// isValidAddress 校验地址是否合法
func isValidAddress(address string) bool {

	//校验一下地址

	//解码，得到25个字节数据
	decodeInfo := base58.Decode(address)

	if len(decodeInfo) != 25 {
		fmt.Println("传入地址无效")
		return false
	}
	//截取前21个payload，截取后4个checksum1
	payload := decodeInfo[:len(decodeInfo)-4]
	checkSum1 := decodeInfo[len(decodeInfo)-4:]
	//对playload计算，得到checkSum2，比较checkSum1和checkSum2,

	checkSum2 := checkSum(payload)

	return bytes.Equal(checkSum1, checkSum2)
}
