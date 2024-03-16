package main

import (
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
	hash1 := sha256.Sum256(pubKey)
	//hash160处理
	//hasher := crypto.RIPEMD160.New()

	hasher := ripemd160.New()

	hasher.Write(hash1[:])
	//公钥哈希，锁定output时就是用这个值
	pubKeyHash := hasher.Sum(nil)
	//拼接version和公钥哈希，得到21个字节
	payload := append([]byte{byte(00)}, pubKeyHash...)

	//生成4字节的校验码
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	//4字节checksum
	checkSum := second[0:4]

	payload = append(payload, checkSum...)

	address := base58.Encode(payload)

	return address
}
